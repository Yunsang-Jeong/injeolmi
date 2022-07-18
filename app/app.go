package app

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xanzy/go-gitlab"
)

type Injeolmi struct {
	gitlabWebhookEvent  string
	gitlabWebhookSecret string
	gitlabToken         string
	gitlabClient        *gitlab.Client
	gitlabWebhookBody   interface{}
	userActionType      string
	userActionOptions   []string
	awsClientConfig     *aws.Config
}

const (
	webhookSecretHeader = "X-Gitlab-Token"
	webhookEventHeader  = "X-Gitlab-Event"
	triggerKeyword      = "injeolmi"
	responseKeyword     = "[Injeolmi]"
	webhookSecretEnvKey = "GITLAB_WEBHOOK_SECRET"
	gitlabTokenEnvKey   = "GITLAB_TOKEN"
	dynamodbTableName   = "ingeolmi-dynamodb"
)

const (
	webhookNoteEvent     = "Note Hook"
	webhookPipelineEvent = "Pipeline Hook"
	webhookJobEvent      = "Job Hook"
)

func Run(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	injeolmi := Injeolmi{}

	// Get gitlab webhook event from header and validate
	allowEventTypeList := []string{webhookNoteEvent, webhookPipelineEvent, webhookJobEvent}
	if err := injeolmi.getGitlabWebookEventAndValidate(req.Headers, allowEventTypeList); err != nil {
		return generateReturn(err.Error(), 400)
	}

	// Get gitlab webhook secret from header and validate
	if err := injeolmi.getGitlabWebookSecretAndValidate(req.Headers); err != nil {
		return generateReturn(err.Error(), 400)
	}

	// Get gitlab token from env
	if err := injeolmi.getGitlabToken(); err != nil {
		return generateReturn(err.Error(), 400)
	}

	// Initiate gitlab client
	if err := injeolmi.setGitlabClient(); err != nil {
		return generateReturn(err.Error(), 400)
	}

	// Initiate aws client
	if err := injeolmi.setAWSClient(); err != nil {
		return generateReturn(err.Error(), 400)
	}

	// Get gitlab webhook body
	if err := injeolmi.getGitlabWebhookBodyAndValidate(req.Body); err != nil {
		return generateReturn(err.Error(), 400)
	}

	// Handle webhook
	if err := injeolmi.handleWebhook(); err != nil {
		return generateReturn(err.Error(), 400)
	}

	return generateReturn("OK", 200)
}

func generateReturn(body interface{}, statusCode int) (events.APIGatewayProxyResponse, error) {
	log.Printf("response body : %s", body)
	log.Printf("response code : %d", statusCode)

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", body), StatusCode: statusCode}, nil
}
