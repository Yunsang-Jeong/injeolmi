package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Yunsang-Jeong/ingeolmi/app"
)

func main() {
	lambda.Start(app.Run)
}
