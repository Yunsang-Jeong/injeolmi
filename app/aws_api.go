package app

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	errors "github.com/pkg/errors"
)

type dynamodb_table_schema map[string]types.AttributeValue

//  {
// 	CommentID     types.AttributeValueMemberS
// 	MRID          types.AttributeValueMemberS
// 	MRIID         types.AttributeValueMemberS
// 	PipelineID    types.AttributeValueMemberS
// 	ActionType    types.AttributeValueMemberS
// 	ActionOptions types.AttributeValueMemberS
// }

func (i *Injeolmi) PutItemToAWSDynamoDB(item dynamodb_table_schema) error {
	cli := dynamodb.NewFromConfig(*i.awsClientConfig)

	// reponse: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#PutItemOutput
	// AWS Docs: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_PutItem.html
	_, err := cli.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(dynamodbTableName),
		Item:      item,
	})
	if err != nil {
		return errors.Wrap(err, "fail to put item to AWS DynamoDB")
	}

	return nil
}

func (i *Injeolmi) GetputItemToAWSDynamoDB(primaryKeyValue int) error {
	cli := dynamodb.NewFromConfig(*i.awsClientConfig)

	_, err := cli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(dynamodbTableName),
		Key: map[string]types.AttributeValue{
			"CommentID": &types.AttributeValueMemberN{Value: strconv.Itoa(primaryKeyValue)},
		},
	})
	if err != nil {
		return errors.Wrap(err, "fail to get item to AWS DynamoDB")
	}

	return nil
}
