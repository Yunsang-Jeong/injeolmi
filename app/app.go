package app

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	errors "github.com/pkg/errors"
)

func Run(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cli := &client{}
	webhook := &gitlabWebhook{}

	// [1] Parse webhook header
	if err := webhook.parseHeader(req.Headers); err != nil {
		return generateReturn(err.Error(), 400)
	}
	log.Printf("success to parse webhook header\n")

	// [2] Init gitlab client
	if err := cli.init(); err != nil {
		return generateReturn(err.Error(), 400)
	}
	log.Printf("success to init gitlab client\n")

	// [3] Validate webhook secret
	if webhook.secret != cli.webhookSecret {
		return generateReturn(errors.New("fail to validate webhook secret"), 403)
	}
	log.Printf("success to validate webhook secret\n")

	// [4] Parse webhook body
	if err := webhook.parseBody(req.Body); err != nil {
		return generateReturn(err.Error(), 400)
	}
	log.Printf("success to parse webhook body\n")

	// [5] Handle webhook
	if err := webhook.handleWebhook(cli); err != nil {
		return generateReturn(err.Error(), 400)
	}
	log.Printf("success to handle webhook\n")

	return generateReturn("OK", 200)
}

func generateReturn(body interface{}, statusCode int) (events.APIGatewayProxyResponse, error) {
	log.Printf("response body : %s", body)
	log.Printf("response code : %d", statusCode)

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("%s", body), StatusCode: statusCode}, nil
}
