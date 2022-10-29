package internal

import (
	"encoding/json"
	"fmt"

	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (c *GitlabEventHandlerConfig) AnalyzeEventBody() error {
	switch c.GitlabEvent.EventType {
	//
	// Analyze "Note Hook" (Comment)
	//
	case "Note Hook":
		var noteType struct {
			ObjectAttributes struct {
				NoteableType string `json:"noteable_type"`
			} `json:"object_attributes"`
		}
		if err := json.Unmarshal([]byte(c.RawEventBody), &noteType); err != nil {
			return errors.Wrap(err, "can't parse comment type on commit from webhook")
		}

		switch noteType.ObjectAttributes.NoteableType {
		case "MergeRequest":
			var e gitlab.MergeCommentEvent
			if err := json.Unmarshal([]byte(c.RawEventBody), &e); err != nil {
				return errors.Wrap(err, "can't parse comment on merge from webhook")
			}

			c.GitlabEvent.EventBody = e

			return nil
		}
	//
	// Analyze "Pipeline Hook"
	//
	case "Pipeline Hook":
		var e gitlab.PipelineEvent
		if err := json.Unmarshal([]byte(c.RawEventBody), &e); err != nil {
			return errors.Wrap(err, "can't parse pipeline event from webhook")
		}
		c.GitlabEvent.EventBody = e

		return nil
	}

	return errors.New(fmt.Sprintf("can't parse %s event from webhook", c.GitlabEvent.EventType))
}
