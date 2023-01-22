package internal

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/Yunsang-Jeong/injeolmi/pkg/awsapi"
	"github.com/Yunsang-Jeong/injeolmi/pkg/vcsapi"
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (c *GitlabEventHandlerConfig) HandleStartComment(event interface{}) error {
	e := event.(gitlab.MergeCommentEvent)

	if len(c.GitlabEvent.ActionOptions) < 1 {
		msg := fmt.Sprintf("[%s] `start` action must have action options!", c.TriggerKeyword)
		_, err := vcsapi.WriteGitlabComment(*c.GitlabClient, e.ProjectID, e.MergeRequest.IID, msg)
		if err != nil {
			return errors.Wrap(err, "fail to handle webhook")
		}

		return nil
	}

	pipelineVariables := &[]*gitlab.PipelineVariable{
		{
			Key:          "CHDIR",
			Value:        c.GitlabEvent.ActionOptions[0],
			VariableType: "env_var",
		},
	}

	pipeline, err := vcsapi.RunGitlabBranchPipeline(*c.GitlabClient, e.ProjectID, e.MergeRequest.SourceBranch, pipelineVariables)
	if err != nil {
		return errors.Wrap(err, "fail to run pipeline")
	}

	commentString := fmt.Sprintf(":zap: **Success to run [pipeliine](%s/) on %s directory**", pipeline.WebURL, strings.Join(c.GitlabEvent.ActionOptions, ", "))
	note, err := vcsapi.WriteGitlabComment(*c.GitlabClient, e.ProjectID, e.MergeRequest.IID, commentString)
	if err != nil {
		return errors.Wrap(err, "fail to write comment")
	}

	item, err := attributevalue.MarshalMap(DynamodbTableSchema{
		PipelineID:       pipeline.ID,
		CommentID:        note.ID,
		MergeRequestsIID: e.MergeRequest.IID,
		ActionType:       "start",
		ActionOptions:    strings.Join(c.GitlabEvent.ActionOptions, ", "),
		CommentString:    commentString,
	})
	if err != nil {
		return errors.Wrap(err, "fail to marshal item to put in dynamodb")
	}

	if err := awsapi.PutItemToAWSDynamodb(*c.AWSConfig, c.DynamodbTableName, item); err != nil {
		return errors.Wrap(err, "fail to put item in dynamodb")
	}

	return nil
}

func (c *GitlabEventHandlerConfig) HandleHelpComment(event interface{}) error {
	e := event.(gitlab.MergeCommentEvent)

	msg := fmt.Sprintf("[%s] help task is not supported yet!", c.TriggerKeyword)

	_, err := vcsapi.WriteGitlabComment(*c.GitlabClient, e.ProjectID, e.MergeRequest.IID, msg)
	if err != nil {
		return errors.Wrap(err, "fail to handle webhook")
	}

	return nil
}
