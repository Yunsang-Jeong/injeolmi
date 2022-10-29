package internal_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/Yunsang-Jeong/ingeolmi/internal"
	"github.com/Yunsang-Jeong/ingeolmi/pkg/awsapi"
	"github.com/Yunsang-Jeong/ingeolmi/pkg/vcsapi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

const (
	gitlabToken         = "glpat-0123456789"
	gitlabWebhookSecret = "SeCrEt!@#"
	awsRegionName       = "ap-northeast-2"
	dynamodbTableName   = "ingeolmi-dynamodb"
	triggerKeyword      = "ingeolmi"
	projectId           = 123456789
	sourceBranch        = "source_branch"
	mergeRequestsIID    = 123456789
	pipelineId          = 123456789
)

var allowedEventTypes = []string{"Note Hook", "Pipeline Hook"}

type gitlabEvent struct {
	projectId        int
	mergeRequestsIID int
	sourceBranch     string
	eventType        string
	comment          string
	pipelineId       int
}

func generateResponse(body interface{}, statusCode int) events.APIGatewayProxyResponse {
	log.Printf("response body : %s", body)
	log.Printf("response code : %d", statusCode)

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", body), StatusCode: statusCode}
}

func (e gitlabEvent) generateNoteHookEvent(t *testing.T) *internal.GitlabEventHandlerConfig {
	body := map[string]interface{}{
		"project_id":  e.projectId,
		"object_kind": "note",
		"merge_request": map[string]interface{}{
			"source_branch": e.sourceBranch,
			"iid":           e.mergeRequestsIID,
		},
		"object_attributes": map[string]interface{}{
			"note":          e.comment,
			"noteable_type": "MergeRequest",
		},
	}
	headers := map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": e.eventType,
	}

	rawbody, _ := json.Marshal(body)

	awsConfig, err := awsapi.GetAWSConfig(awsRegionName)
	assert.IsType(t, nil, err)

	gitlabClient, err := vcsapi.GetGitlabClient(gitlabToken)
	assert.IsType(t, nil, err)

	return &internal.GitlabEventHandlerConfig{
		GitlabWebhookSecret: gitlabWebhookSecret,
		DynamodbTableName:   dynamodbTableName,
		TriggerKeyword:      triggerKeyword,
		AllowedEventTypes:   allowedEventTypes,
		AWSConfig:           awsConfig,
		GitlabClient:        gitlabClient,
		RawEventHeaders:     headers,
		RawEventBody:        string(rawbody[:]),
	}
}

func (e gitlabEvent) generateValidateSuccessPipelineEvents(t *testing.T) *internal.GitlabEventHandlerConfig {
	body := map[string]interface{}{
		"project": map[string]interface{}{
			"id": e.projectId,
		},
		"object_kind": "pipeline",
		"merge_request": map[string]interface{}{
			"source_branch": e.sourceBranch,
			"iid":           e.mergeRequestsIID,
		},
		"object_attributes": map[string]interface{}{
			"id":     e.pipelineId,
			"status": "running",
			"source": "api",
		},
		"builds": []map[string]interface{}{
			{
				"id":     2,
				"stage":  "plan",
				"name":   "terraform plan",
				"status": "created",
				"when":   "on_success",
				"manual": false,
			},
			{
				"id":     3,
				"stage":  "apply",
				"name":   "terraform apply",
				"status": "created",
				"when":   "manual",
				"manual": true,
			},
			{
				"id":     1,
				"stage":  "validate",
				"name":   "terraform validate",
				"status": "success",
				"when":   "on_success",
				"manual": false,
			},
		},
	}

	headers := map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": e.eventType,
	}

	rawbody, _ := json.Marshal(body)

	awsConfig, err := awsapi.GetAWSConfig(awsRegionName)
	assert.IsType(t, nil, err)

	gitlabClient, err := vcsapi.GetGitlabClient(gitlabToken)
	assert.IsType(t, nil, err)

	return &internal.GitlabEventHandlerConfig{
		GitlabWebhookSecret: gitlabWebhookSecret,
		DynamodbTableName:   dynamodbTableName,
		TriggerKeyword:      triggerKeyword,
		AllowedEventTypes:   allowedEventTypes,
		AWSConfig:           awsConfig,
		GitlabClient:        gitlabClient,
		RawEventHeaders:     headers,
		RawEventBody:        string(rawbody[:]),
	}
}

func (e gitlabEvent) generatePlanSuccessPipelineEvents(t *testing.T) *internal.GitlabEventHandlerConfig {
	body := map[string]interface{}{
		"project": map[string]interface{}{
			"id": e.projectId,
		},
		"object_kind": "pipeline",
		"merge_request": map[string]interface{}{
			"source_branch": e.sourceBranch,
			"iid":           e.mergeRequestsIID,
		},
		"object_attributes": map[string]interface{}{
			"id":     e.pipelineId,
			"status": "running",
			"source": "api",
		},
		"builds": []map[string]interface{}{
			{
				"id":     2,
				"stage":  "plan",
				"name":   "terraform plan",
				"status": "success",
				"when":   "on_success",
				"manual": false,
			},
			{
				"id":     3,
				"stage":  "apply",
				"name":   "terraform apply",
				"status": "created",
				"when":   "manual",
				"manual": true,
			},
			{
				"id":     1,
				"stage":  "validate",
				"name":   "terraform validate",
				"status": "success",
				"when":   "on_success",
				"manual": false,
			},
		},
	}

	headers := map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": e.eventType,
	}

	rawbody, _ := json.Marshal(body)

	awsConfig, err := awsapi.GetAWSConfig(awsRegionName)
	assert.IsType(t, nil, err)

	gitlabClient, err := vcsapi.GetGitlabClient(gitlabToken)
	assert.IsType(t, nil, err)

	return &internal.GitlabEventHandlerConfig{
		GitlabWebhookSecret: gitlabWebhookSecret,
		DynamodbTableName:   dynamodbTableName,
		TriggerKeyword:      triggerKeyword,
		AllowedEventTypes:   allowedEventTypes,
		AWSConfig:           awsConfig,
		GitlabClient:        gitlabClient,
		RawEventHeaders:     headers,
		RawEventBody:        string(rawbody[:]),
	}
}

func (e gitlabEvent) generateApplySuccessPipelineEvents(t *testing.T) *internal.GitlabEventHandlerConfig {
	body := map[string]interface{}{
		"project": map[string]interface{}{
			"id": e.projectId,
		},
		"object_kind": "pipeline",
		"merge_request": map[string]interface{}{
			"source_branch": e.sourceBranch,
			"iid":           e.mergeRequestsIID,
		},
		"object_attributes": map[string]interface{}{
			"id":     e.pipelineId,
			"status": "running",
			"source": "api",
		},
		"builds": []map[string]interface{}{
			{
				"id":     2,
				"stage":  "plan",
				"name":   "terraform plan",
				"status": "success",
				"when":   "on_success",
				"manual": false,
			},
			{
				"id":     3,
				"stage":  "apply",
				"name":   "terraform apply",
				"status": "success",
				"when":   "manual",
				"manual": true,
			},
			{
				"id":     1,
				"stage":  "validate",
				"name":   "terraform validate",
				"status": "success",
				"when":   "on_success",
				"manual": false,
			},
		},
	}

	headers := map[string]string{
		"X-Gitlab-Token": gitlabWebhookSecret,
		"X-Gitlab-Event": e.eventType,
	}

	rawbody, _ := json.Marshal(body)

	awsConfig, err := awsapi.GetAWSConfig(awsRegionName)
	assert.IsType(t, nil, err)

	gitlabClient, err := vcsapi.GetGitlabClient(gitlabToken)
	assert.IsType(t, nil, err)

	return &internal.GitlabEventHandlerConfig{
		GitlabWebhookSecret: gitlabWebhookSecret,
		DynamodbTableName:   dynamodbTableName,
		TriggerKeyword:      triggerKeyword,
		AllowedEventTypes:   allowedEventTypes,
		AWSConfig:           awsConfig,
		GitlabClient:        gitlabClient,
		RawEventHeaders:     headers,
		RawEventBody:        string(rawbody[:]),
	}
}
