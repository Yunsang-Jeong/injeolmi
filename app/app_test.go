package app_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Yunsang-Jeong/ingeolmi/app"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type testSet struct {
	request events.APIGatewayProxyRequest
	expect  events.APIGatewayProxyResponse
	err     error
}

const gitlabToken = "glpat-01234567890123456789"
const gitlabWebhookSecret = "SeCrEt!@#"
const gitlabEventNoteHook = "Note Hook"
const gitlabEventPipeline = "Pipeline Hook"
const projectId = 31984415
const mergeRequestsIID = 12
const sourceBranch = "terraform-test"

func TestRun_task_help(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", gitlabToken)
	os.Setenv("GITLAB_WEBHOOK_SECRET", gitlabWebhookSecret)

	var mockRequestHeader = map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": gitlabEventNoteHook,
	}

	var mockRequestBody = map[string]interface{}{
		"project_id":  projectId,
		"object_kind": "note",
		"merge_request": map[string]interface{}{
			"source_branch": sourceBranch,
			"iid":           mergeRequestsIID,
		},
		"object_attributes": map[string]string{
			"note":          "injeolmi help",
			"noteable_type": "MergeRequest",
		},
	}
	marshaledMockRequestBody, _ := json.Marshal(mockRequestBody)

	testSet := testSet{
		request: events.APIGatewayProxyRequest{
			Headers: mockRequestHeader,
			Body:    string(marshaledMockRequestBody[:]),
		},
		expect: events.APIGatewayProxyResponse{
			Body:       "fail to handle webhook: fail to write comment: POST https://gitlab.com/api/v4/projects/31984415/merge_requests/12/notes: 401 {message: 401 Unauthorized}",
			StatusCode: 400,
		},
		err: nil,
	}

	response, err := app.Run(testSet.request)
	assert.IsType(t, testSet.err, err)
	assert.Equal(t, testSet.expect.Body, response.Body)
	assert.Equal(t, testSet.expect.StatusCode, response.StatusCode)
}

func TestRun_task_start(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", gitlabToken)
	os.Setenv("GITLAB_WEBHOOK_SECRET", gitlabWebhookSecret)

	var mockRequestHeader = map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": gitlabEventNoteHook,
	}

	var mockRequestBody = map[string]interface{}{
		"project_id":  projectId,
		"object_kind": "note",
		"merge_request": map[string]interface{}{
			"source_branch": sourceBranch,
			"iid":           mergeRequestsIID,
		},
		"object_attributes": map[string]string{
			"note":          "injeolmi start /app/dev",
			"noteable_type": "MergeRequest",
		},
	}
	marshaledMockRequestBody, _ := json.Marshal(mockRequestBody)

	testSet := testSet{
		request: events.APIGatewayProxyRequest{
			Headers: mockRequestHeader,
			Body:    string(marshaledMockRequestBody[:]),
		},
		expect: events.APIGatewayProxyResponse{
			Body:       "fail to handle webhook in start: fail to run branch pipeline: POST https://gitlab.com/api/v4/projects/31984415/pipeline: 401 {message: 401 Unauthorized}",
			StatusCode: 400,
		},
		err: nil,
	}

	response, err := app.Run(testSet.request)
	assert.IsType(t, testSet.err, err)
	assert.Equal(t, testSet.expect.Body, response.Body)
	assert.Equal(t, testSet.expect.StatusCode, response.StatusCode)
}

func TestRun_updateComment(t *testing.T) {
	os.Setenv("GITLAB_TOKEN", gitlabToken)
	os.Setenv("GITLAB_WEBHOOK_SECRET", gitlabWebhookSecret)

	var mockRequestHeader = map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": gitlabEventPipeline,
	}

	var mockRequestBody = map[string]interface{}{
		"project_id":  projectId,
		"object_kind": "pipeline",
		"merge_request": map[string]interface{}{
			"source_branch": sourceBranch,
			"iid":           mergeRequestsIID,
		},
		"object_attributes": map[string]string{
			"source":        "api",
			"noteable_type": "MergeRequest",
		},
	}
	marshaledMockRequestBody, _ := json.Marshal(mockRequestBody)

	testSet := testSet{
		request: events.APIGatewayProxyRequest{
			Headers: mockRequestHeader,
			Body:    string(marshaledMockRequestBody[:]),
		},
		expect: events.APIGatewayProxyResponse{
			Body:       "fail to handle webhook in start: fail to run branch pipeline: POST https://gitlab.com/api/v4/projects/31984415/pipeline: 401 {message: 401 Unauthorized}",
			StatusCode: 400,
		},
		err: nil,
	}

	response, err := app.Run(testSet.request)
	assert.IsType(t, testSet.err, err)
	assert.Equal(t, testSet.expect.Body, response.Body)
	assert.Equal(t, testSet.expect.StatusCode, response.StatusCode)
}
