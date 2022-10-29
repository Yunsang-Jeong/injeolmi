package main

import (
	"fmt"
	"log"

	"github.com/Yunsang-Jeong/ingeolmi/internal"
	"github.com/Yunsang-Jeong/ingeolmi/pkg/awsapi"
	"github.com/Yunsang-Jeong/ingeolmi/pkg/vcsapi"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/caarlos0/env"
	"github.com/xanzy/go-gitlab"
)

var awsConfig *aws.Config
var gitlabClient *gitlab.Client

type lambdaEnv struct {
	GitlabToken         string   `env:"GITLAB_TOKEN,required"`
	GitlabWebhookSecret string   `env:"GITLAB_WEBHOOK_SECRET,required"`
	AWSRegionName       string   `env:"AWS_REGION_NAME" envDefault:"ap-northeast-2"`
	DynamodbTableName   string   `env:"DYNAMODB_TABLE_NAME,required"`
	TriggerKeyword      string   `env:"TRIGGER_KEYWORD" envDefault:"injeolmi"`
	AllowedEventTypes   []string `env:"ALLOWED_EVENT_TYPES" envSeparator:"," envDefault:"Note Hook,Pipeline Hook"`
}

func main() {
	var err error

	e := lambdaEnv{}
	if err := env.Parse(&e); err != nil {
		panic(err)
	}

	awsConfig, err = awsapi.GetAWSConfig(e.AWSRegionName)
	if err != nil {
		panic(err)
	}

	gitlabClient, err = vcsapi.GetGitlabClient(e.GitlabToken)
	if err != nil {
		panic(err)
	}

	lambda.Start(e.LambdaStart)
}

func (e lambdaEnv) LambdaStart(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return internal.Run(
		&internal.GitlabEventHandlerConfig{
			/*
				Get from ENV
			*/
			GitlabWebhookSecret: e.GitlabWebhookSecret,
			DynamodbTableName:   e.DynamodbTableName,
			TriggerKeyword:      e.TriggerKeyword,
			AllowedEventTypes:   e.AllowedEventTypes,
			/*
				Get from global variables
			*/
			AWSConfig:    awsConfig,
			GitlabClient: gitlabClient,
			/*
				Get from event
			*/
			RawEventHeaders: req.Headers,
			RawEventBody:    req.Body,
		},
		generateResponse,
	)
}

func generateResponse(body interface{}, statusCode int) events.APIGatewayProxyResponse {
	log.Printf("response body : %s", body)
	log.Printf("response code : %d", statusCode)

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", body), StatusCode: statusCode}
}
