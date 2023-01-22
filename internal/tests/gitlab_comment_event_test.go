package internal_test

import (
	"fmt"
	"testing"

	"github.com/Yunsang-Jeong/injeolmi/internal"
	"github.com/stretchr/testify/assert"
)

func Test_KeywordComment(t *testing.T) {
	e := gitlabEvent{
		projectId:        projectId,
		mergeRequestsIID: mergeRequestsIID,
		sourceBranch:     sourceBranch,
		eventType:        "Note Hook",
		comment:          triggerKeyword,
	}

	c := e.generateNoteHookEvent(t)
	r, err := internal.Run(c, generateResponse)
	assert.IsType(t, nil, err)
	assert.Equal(t, 200, r.StatusCode)
}

func Test_HelpComment(t *testing.T) {
	e := gitlabEvent{
		projectId:        projectId,
		mergeRequestsIID: mergeRequestsIID,
		sourceBranch:     sourceBranch,
		eventType:        "Note Hook",
		comment:          fmt.Sprintf("%s help", triggerKeyword),
	}

	c := e.generateNoteHookEvent(t)
	r, err := internal.Run(c, generateResponse)
	assert.IsType(t, nil, err)
	assert.Equal(t, 200, r.StatusCode)
}

func Test_StartComment(t *testing.T) {
	e := gitlabEvent{
		projectId:        projectId,
		mergeRequestsIID: mergeRequestsIID,
		sourceBranch:     sourceBranch,
		eventType:        "Note Hook",
		comment:          fmt.Sprintf("%s start app/prd", triggerKeyword),
	}

	c := e.generateNoteHookEvent(t)
	r, err := internal.Run(c, generateResponse)
	assert.IsType(t, nil, err)
	assert.Equal(t, 200, r.StatusCode)
}
