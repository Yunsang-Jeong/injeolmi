package internal

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	errors "github.com/pkg/errors"

	"github.com/xanzy/go-gitlab"
)

type GitlabEventHandlerConfig struct {
	GitlabWebhookSecret string
	DynamodbTableName   string
	AWSConfig           *aws.Config
	GitlabClient        *gitlab.Client
	TriggerKeyword      string
	AllowedEventTypes   []string
	RawEventHeaders     map[string]string
	RawEventBody        string
	GitlabEvent         GitlabEvent
}

type GitlabEvent struct {
	EventType     string
	WebhookSecret string
	EventBody     interface{}
	ActionType    string
	ActionOptions []string
}

func (c *GitlabEventHandlerConfig) HandleEvent() error {
	if c.GitlabEvent.EventBody == nil {
		return errors.New("[HandleEvent] Eventbody is empty!")
	}

	switch e := c.GitlabEvent.EventBody.(type) {
	/*

		Merge Comment Event

	*/
	case gitlab.MergeCommentEvent:
		if err := c.AnalyzeComment(e.ObjectAttributes.Note); err != nil {
			return err
		}

		switch c.GitlabEvent.ActionType {
		case "start":
			if err := c.HandleStartComment(e); err != nil {
				return errors.Wrap(err, "fail to handle webhook in start")
			}

		case "help":
			if err := c.HandleHelpComment(e); err != nil {
				return errors.Wrap(err, "fail to handle webhook in help")
			}
		}
	/*

		Pipeline Event

	*/
	case gitlab.PipelineEvent:
		// status: pending --> running --> success(passed)

		for _, build := range e.Builds {
			log.Printf("[PipelineEvent] %s/%s is %s", build.Stage, build.Name, build.Status)
		}

		if err := c.HandlePipeline(e); err != nil {
			return errors.Wrap(err, "fail to handle webhook in pipelineevent")
		}

	default:
		log.Printf("%T event is not supported", e)
	}

	return nil
}
