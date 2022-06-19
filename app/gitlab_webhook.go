package app

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type gitlabWebhook struct {
	secret string
	event  string
	body   interface{}
}

func (h *gitlabWebhook) parseHeader(header map[string]string) error {
	const webhookSecretHeader = "X-Gitlab-Token"
	const webhookEventHeader = "X-Gitlab-Event"

	if secret, exist := header[webhookSecretHeader]; exist {
		h.secret = secret
	} else {
		return errors.New("can't parse secret from webhook")
	}

	if event, exist := header[webhookEventHeader]; exist {
		h.event = event
	} else {
		return errors.New("can't parse event from webhook")
	}

	return nil
}

func (h *gitlabWebhook) parseBody(body string) error {
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
		h.body = e
	}

	return nil
}

func (h *gitlabWebhook) handleWebhook(c *client) error {
	const fmtPipelineSuccessString = "[Injeolmi] Success to run [pipeliine](%s/)"
	const helpString = "[Injeolmi] help task is not supported yet!"

	switch body := h.body.(type) {
	case gitlab.MergeCommentEvent:
		//
		// handle MergeCommentEvent
		//
		comment := body.ObjectAttributes.Note
		commentSlice := strings.Split(comment, " ")

		if commentSlice[0] != "injeolmi" {
			return errors.New("fail to handle webhook. comment does not start with injeolmi.")
		}

		if len(commentSlice) < 2 {
			return errors.New("fail to handle webhook. comment need something more!")
		}

		switch commentSlice[1] {
		case "start":
			//
			// injeolmi start %DIRECTORY%
			//
			variables := &[]*gitlab.PipelineVariable{
				{Key: "DIRECTORY", Value: commentSlice[2], VariableType: "env_var"},
			}

			if err := c.RunBranchPipeline(body.ProjectID, body.MergeRequest.SourceBranch, variables); err != nil {
				return errors.Wrap(err, "fail to handle webhook in start")
			}

			pipeline := c.response["pipeline"].(*gitlab.Pipeline)

			if err := c.WriteComment(body.ProjectID, body.MergeRequest.IID, fmt.Sprintf(fmtPipelineSuccessString, pipeline.WebURL)); err != nil {
				return errors.Wrap(err, "fail to handle webhook")
			}

		case "help":
			//
			// injeolmi help
			//
			if err := c.WriteComment(body.ProjectID, body.MergeRequest.IID, helpString); err != nil {
				return errors.Wrap(err, "fail to handle webhook")
			}
		}

	default:
		log.Printf("%T event is not supported", body)
	}

	return nil
}
