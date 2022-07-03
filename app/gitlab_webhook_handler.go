package app

import (
	"fmt"
	"log"

	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

const (
	fmtPipelineSuccessString = "[Injeolmi] Success to run [pipeliine](%s/)"
	helpString               = "[Injeolmi] help task is not supported yet!"
)

func (i *Injeolmi) handleWebhook() error {
	switch body := i.gitlabWebhookBody.(type) {
	case gitlab.MergeCommentEvent:
		if err := i.parseUserActionFromMRComment(body.ObjectAttributes.Note); err != nil {
			return err
		}

		switch i.userActionType {
		case "start":
			i.handleStartAction(body.ProjectID, body.MergeRequest.SourceBranch, body.MergeRequest.IID)

		case "help":
			i.handleHelpAction(body.ProjectID, body.MergeRequest.IID)
		}

	case gitlab.PipelineEvent:
		i.handlePipelineEvent()

	case gitlab.JobStats:
		i.handleJobEvent()

	default:
		log.Printf("%T event is not supported", body)
	}

	return nil
}

func (i *Injeolmi) handleStartAction(projectId int, sourceBranch string, mergeRequestIID int) error {
	if i.userActionOptions == nil {
		return errors.New("`start` action must have action options")
	}

	variables := &[]*gitlab.PipelineVariable{
		{
			Key:          "DIR",
			Value:        i.userActionOptions[0],
			VariableType: "env_var",
		},
	}

	if err := i.runBranchPipeline(projectId, sourceBranch, variables); err != nil {
		return errors.Wrap(err, "fail to handle webhook in start")
	}

	pipeline := i.gitlabClientResponse["pipeline"].(*gitlab.Pipeline)
	// pipeline.ID

	if err := i.writeComment(projectId, mergeRequestIID, fmt.Sprintf(fmtPipelineSuccessString, pipeline.WebURL)); err != nil {
		return errors.Wrap(err, "fail to handle webhook")
	}

	// item := dynamodb_table_attributes{
	// 	CommentID: body.ObjectAttributes.ID,
	// 	MRID: body.MergeRequest.ID,
	// 	MRIID: body.MergeRequest.IID,
	// 	PipelineID:
	// 	ActionType: h.ActionType,
	// 	ActionOptions: h.AcActionOptions,
	// }

	return nil
}

func (i *Injeolmi) handleHelpAction(projectId int, mergeRequestIID int) error {
	if err := i.writeComment(projectId, mergeRequestIID, helpString); err != nil {
		return errors.Wrap(err, "fail to handle webhook")
	}

	return nil
}

func (i *Injeolmi) handlePipelineEvent() error {
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

func (i *Injeolmi) handleJobEvent() error {
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
