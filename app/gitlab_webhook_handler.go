package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

const (
	fmtPipelineSuccessString = "[Injeolmi] Success to run [pipeliine](%s/)"
	helpString               = "[Injeolmi] help task is not supported yet!"
)

func (i *Injeolmi) handleWebhook() error {
	switch e := i.gitlabWebhookBody.(type) {
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

	case gitlab.PipelineEvent:
		if err := i.handlePipelineEvent(e); err != nil {
			return errors.Wrap(err, "fail to handle webhook in pipelineevent")
		}

	case gitlab.JobStats:
		if err := i.handleJobEvent(e); err != nil {
			return errors.Wrap(err, "fail to handle webhook in jobstats")
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
	log.Println(i.userActionOptions[0])
	pipeline, err := i.runBranchPipeline(e.ProjectID, e.MergeRequest.SourceBranch, pipelineVariables)
	if err != nil {
		return errors.Wrap(err, "fail to run pipeline")
	}

	//
	// Write comment on mr
	//
	commentString := fmt.Sprintf(fmtPipelineSuccessString, pipeline.WebURL)
	note, err := i.writeComment(e.ProjectID, e.MergeRequest.IID, commentString)
	if err != nil {
		return errors.Wrap(err, "fail to write comment")
	}

	//
	// Save data to AWS DynamoDB
	//
	item := dynamodb_table_schema{
		"ProjectID":        &types.AttributeValueMemberN{Value: strconv.Itoa(e.ProjectID)},
		"CommentID":        &types.AttributeValueMemberN{Value: strconv.Itoa(note.ID)},
		"MergeRequestsID":  &types.AttributeValueMemberN{Value: strconv.Itoa(e.MergeRequest.ID)},
		"MergeRequestsIID": &types.AttributeValueMemberN{Value: strconv.Itoa(e.MergeRequest.IID)},
		"PipelineID":       &types.AttributeValueMemberN{Value: strconv.Itoa(pipeline.ID)},
		"ActionType":       &types.AttributeValueMemberS{Value: "start"},
		"ActionOptions":    &types.AttributeValueMemberS{Value: strings.Join(i.userActionOptions, ", ")},
	}
	if err := i.PutItemToAWSDynamoDB(item); err != nil {
		return errors.Wrap(err, "fail to put item to AWS DynamoDB")
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
	// pipelineId := body.ObjectAttributes.ID
	// pipelineVariables := body.ObjectAttributes.Variables

	// dir := ""
	// for _, variable := range pipelineVariables {
	// 	if variable.Key == "DIR" {
	// 		dir = variable.Value
	// 	}
	// }
	return nil
}

func (i *Injeolmi) handleJobEvent(e gitlab.JobStats) error {
	//
	// handle PipelineEvent
	//
	// pipelineId := body.ObjectAttributes.ID
	// pipelineVariables := body.ObjectAttributes.Variables

	// dir := ""
	// for _, variable := range pipelineVariables {
	// 	if variable.Key == "DIR" {
	// 		dir = variable.Value
	// 	}
	// }
	return nil
}
