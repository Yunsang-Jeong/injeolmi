package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

const (
	fmtPipelineSuccessString = ":zap: **Success to run [pipeliine](%s/) on %s directory**"
	helpString               = "[Injeolmi] help task is not supported yet!"
)

func (i *Injeolmi) handleWebhook() error {
	switch e := i.gitlabWebhookBody.(type) {
	//
	// handle merge comment event
	//
	case gitlab.MergeCommentEvent:
		if err := i.parseUserActionFromMRComment(e.ObjectAttributes.Note); err != nil {
			return err
		}

		switch i.userActionType {
		case responseKeyword:
			return nil

		case "start":
			if err := i.handleStartAction(e); err != nil {
				return errors.Wrap(err, "fail to handle webhook in start")
			}

		case "help":
			if err := i.handleHelpAction(e); err != nil {
				return errors.Wrap(err, "fail to handle webhook in help")
			}
		}

	//
	// handle pipeline event
	//
	case gitlab.PipelineEvent:
		// status: pending --> running --> success(passed)

		for _, build := range e.Builds {
			log.Printf("[PipelineEvent] %s/%s is %s", build.Stage, build.Name, build.Status)
		}

		if err := i.handlePipelineEvent(e); err != nil {
			return errors.Wrap(err, "fail to handle webhook in pipelineevent")
		}

	default:
		log.Printf("%T event is not supported", e)
	}

	return nil
}

func (i *Injeolmi) handleStartAction(e gitlab.MergeCommentEvent) error {
	if i.userActionOptions == nil {
		return errors.New("`start` action must have action options")
	}

	//
	// Run pipeline
	//
	pipelineVariables := &[]*gitlab.PipelineVariable{
		{
			Key:          "CHDIR",
			Value:        i.userActionOptions[0],
			VariableType: "env_var",
		},
	}

	pipeline, err := i.runBranchPipeline(e.ProjectID, e.MergeRequest.SourceBranch, pipelineVariables)
	if err != nil {
		return errors.Wrap(err, "fail to run pipeline")
	}

	//
	// Write comment on mr
	//
	commentString := fmt.Sprintf(fmtPipelineSuccessString, pipeline.WebURL, strings.Join(i.userActionOptions, ", "))
	note, err := i.writeComment(e.ProjectID, e.MergeRequest.IID, commentString)
	if err != nil {
		return errors.Wrap(err, "fail to write comment")
	}

	//
	// Save data to AWS DynamoDB
	//
	item := DynamodbTableFields{
		PipelineID:       pipeline.ID,
		CommentID:        note.ID,
		MergeRequestsIID: e.MergeRequest.IID,
		ActionType:       "start",
		ActionOptions:    strings.Join(i.userActionOptions, ", "),
		CommentString:    commentString,
	}
	marshaledItem, err := item.marshalDynamodbAttributeValue()
	if err != nil {
		return errors.Wrap(err, "fail to marshal item to put in dynamodb")
	}

	if err := i.PutItemToAWSDynamodb(marshaledItem); err != nil {
		return errors.Wrap(err, "fail to put item in dynamodb")
	}

	return nil
}

func (i *Injeolmi) handleHelpAction(e gitlab.MergeCommentEvent) error {
	//
	// Write comment on mr
	//
	_, err := i.writeComment(e.ProjectID, e.MergeRequest.IID, helpString)
	if err != nil {
		return errors.Wrap(err, "fail to handle webhook")
	}

	return nil
}

func (i *Injeolmi) handlePipelineEvent(e gitlab.PipelineEvent) error {
	//
	// handle PipelineEvent
	//

	item := DynamodbTableFields{
		PipelineID: e.ObjectAttributes.ID,
	}
	marshaledItem, err := item.marshalDynamodbAttributeValue()
	if err != nil {
		return errors.Wrap(err, "fail to marshal item to put in dynamodb")
	}

	returnedItem, err := i.GetItemToAWSDynamodb(marshaledItem)
	if err != nil {
		return errors.Wrap(err, "fail to get item to AWS DynamoDB")
	}

	unmarshaledReturnedItem := make(map[string]string)

	for itemKey, itemValue := range *returnedItem {
		var v string
		if err := attributevalue.Unmarshal(itemValue, &v); err != nil {
			return errors.Wrap(err, "can't parse unmarshal")
		}

		unmarshaledReturnedItem[itemKey] = v
	}

	mergeRequestsIIDm, _ := strconv.Atoi(unmarshaledReturnedItem["MergeRequestsIID"])
	commentID, _ := strconv.Atoi(unmarshaledReturnedItem["CommentID"])
	commentString := unmarshaledReturnedItem["CommentString"]

	for _, build := range e.Builds {
		prefix_emoji := ":white_large_square:"

		switch build.Status {
		case "success":
			prefix_emoji = ":white_check_mark:"
		case "running":
			prefix_emoji = ":runner:"
		}

		c := fmt.Sprintf("- %s - %s is %s : %s", prefix_emoji, build.Stage, build.Name, build.Status)
		commentString = fmt.Sprintf("%s\n%s", commentString, c)
	}

	i.updateComment(
		e.Project.ID,
		mergeRequestsIIDm,
		commentID,
		commentString,
	)

	return nil
}
