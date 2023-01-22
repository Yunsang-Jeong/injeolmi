package internal_test

import (
	"testing"

	"github.com/Yunsang-Jeong/injeolmi/internal"
	"github.com/stretchr/testify/assert"
)

func Test_PipelineTest(t *testing.T) {
	e := gitlabEvent{
		projectId:        projectId,
		mergeRequestsIID: mergeRequestsIID,
		sourceBranch:     sourceBranch,
		eventType:        "Pipeline Hook",
		pipelineId:       pipelineId,
	}

	c1 := e.generateValidateSuccessPipelineEvents(t)
	r1, err1 := internal.Run(c1, generateResponse)
	assert.IsType(t, nil, err1)
	assert.Equal(t, 200, r1.StatusCode)

	c2 := e.generatePlanSuccessPipelineEvents(t)
	r2, err2 := internal.Run(c2, generateResponse)
	assert.IsType(t, nil, err2)
	assert.Equal(t, 200, r2.StatusCode)

	c3 := e.generateApplySuccessPipelineEvents(t)
	r3, err3 := internal.Run(c3, generateResponse)
	assert.IsType(t, nil, err3)
	assert.Equal(t, 200, r3.StatusCode)
}
