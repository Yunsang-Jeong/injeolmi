package internal

import (
	"github.com/aws/aws-lambda-go/events"
)

type responseGenerator func(interface{}, int) events.APIGatewayProxyResponse

func Run(c *GitlabEventHandlerConfig, r responseGenerator) (events.APIGatewayProxyResponse, error) {

	if err := c.AnalyzeEventHeaders(); err != nil {
		return r(err, 400), nil
	}

	if err := c.AnalyzeEventBody(); err != nil {
		return r(err, 400), nil
	}

	if err := c.HandleEvent(); err != nil {
		return r(err, 400), nil
	}

	return r("OK", 200), nil
}
