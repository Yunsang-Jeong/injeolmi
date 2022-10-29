package internal

import (
	"fmt"

	errors "github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func (c *GitlabEventHandlerConfig) AnalyzeEventHeaders() error {
	var exist bool

	// Get webhook event type from header
	c.GitlabEvent.EventType, exist = c.RawEventHeaders["X-Gitlab-Event"]
	if !exist {
		return errors.New("can't parse 'X-Gitlab-Event' from webhook")
	}

	// Get webhook secret from header
	c.GitlabEvent.WebhookSecret, exist = c.RawEventHeaders["X-Gitlab-Token"]
	if !exist {
		return errors.New("can't parse 'X-Gitlab-Token' from webhook")
	}

	// Validate event type
	if !slices.Contains(c.AllowedEventTypes, c.GitlabEvent.EventType) {
		return errors.New(fmt.Sprintf("%s is not allowed event type", c.GitlabEvent.EventType))
	}

	// Validate webhook secret
	if c.GitlabEvent.WebhookSecret != c.GitlabWebhookSecret {
		return errors.New("webhook do not match")
	}

	return nil
}
