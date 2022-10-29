package internal

import (
	"strings"
)

func (c *GitlabEventHandlerConfig) AnalyzeComment(comment string) error {
	commentSlice := strings.Split(comment, " ")

	if c.TriggerKeyword != commentSlice[0] {
		return nil
	}

	if len(commentSlice) == 1 {
		c.GitlabEvent.ActionType = "help"
	} else if len(commentSlice) > 1 {
		c.GitlabEvent.ActionType = commentSlice[1]
		c.GitlabEvent.ActionOptions = commentSlice[2:]
	}

	return nil
}
