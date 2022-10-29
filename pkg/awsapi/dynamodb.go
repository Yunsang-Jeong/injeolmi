package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	errors "github.com/pkg/errors"
)

func PutItemToAWSDynamodb(cfg aws.Config, dynamodbTableName string, item map[string]types.AttributeValue) error {
	cli := dynamodb.NewFromConfig(cfg)

	// reponse: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#PutItemOutput
	// AWS Docs: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_PutItem.html
	_, err := cli.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(dynamodbTableName),
			Item:      item,
		},
	)
	if err != nil {
		return errors.Wrap(err, "fail to put item to AWS DynamoDB")
	}

	return nil
}

func GetItemToAWSDynamodb(cfg aws.Config, dynamodbTableName string, item map[string]types.AttributeValue) (*map[string]types.AttributeValue, error) {
	cli := dynamodb.NewFromConfig(cfg)

	resp, err := cli.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			TableName: aws.String(dynamodbTableName),
			Key:       item,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get item to AWS DynamoDB")
	}

	return &resp.Item, nil
}
