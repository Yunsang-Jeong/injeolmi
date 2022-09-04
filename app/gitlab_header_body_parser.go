package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (i *Injeolmi) getGitlabWebookEventAndValidate(headers map[string]string, allowEventTypeList []string) error {
	// Get webhook event from header
	if event, exist := headers[webhookEventHeader]; exist {
		i.gitlabWebhookEvent = event
	} else {
		return errors.New(fmt.Sprintf("can't parse %s from webhook", webhookEventHeader))
	}

	// Validate event
	for _, eventType := range allowEventTypeList {
		if i.gitlabWebhookEvent == eventType {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("%s is not supported webhook type", i.gitlabWebhookEvent))
}

func (i *Injeolmi) getGitlabWebookSecretAndValidate(headers map[string]string) error {
	// Get webhook secret from header
	if secret, exist := headers[webhookSecretHeader]; exist {
		i.gitlabWebhookSecret = secret
	} else {
		return errors.New("can't parse secret from webhook")
	}

	// Get gitlab webhook secret from env
	gitlabWebhookSecretFromEnv := os.Getenv(webhookSecretEnvKey)
	if gitlabWebhookSecretFromEnv == "" {
		return errors.New("can't load gitlab webhook secret from env")
	}

	// Validate webhook secret
	if i.gitlabWebhookSecret != gitlabWebhookSecretFromEnv {
		return errors.New("fail to validate gitlab webhook secret")
	}

	return nil
}

func (i *Injeolmi) getGitlabWebhookBodyAndValidate(body string) error {
	switch i.gitlabWebhookEvent {

	//
	// parse note event
	//
	case webhookNoteEvent:
		var noteType struct {
			ObjectAttributes struct {
				NoteableType string `json:"noteable_type"`
			} `json:"object_attributes"`
		}
		if err := json.Unmarshal([]byte(body), &noteType); err != nil {
			return errors.Wrap(err, "can't parse comment type on commit from webhook")
		}

		switch noteType.ObjectAttributes.NoteableType {
		case "MergeRequest":
			var e gitlab.MergeCommentEvent
			if err := json.Unmarshal([]byte(body), &e); err != nil {
				return errors.Wrap(err, "can't parse comment on merge from webhook")
			}
			i.gitlabWebhookBody = e
			return nil
		}

	//
	// parse pipeline event
	//
	case webhookPipelineEvent:
		var e gitlab.PipelineEvent
		if err := json.Unmarshal([]byte(body), &e); err != nil {
			return errors.Wrap(err, "can't parse pipeline event from webhook")
		}
		i.gitlabWebhookBody = e
		return nil
	}

	return errors.New(fmt.Sprintf("can't parse %s event from webhook", i.gitlabWebhookEvent))
}

func (i *Injeolmi) parseUserActionFromMRComment(comment string) error {
	commentSlice := strings.Split(comment, " ")

	if responseKeyword == commentSlice[0] {
		return nil
	}

	if triggerKeyword != commentSlice[0] {
		return errors.New("Comments must start with trigger keyword")
	}

	if len(commentSlice) < 2 {
		return errors.New("Comments must contain action keyword")
	}
	i.userActionType = commentSlice[1]

	if 2 < len(commentSlice) {
		i.userActionOptions = commentSlice[2:]
	}

	return nil
}
