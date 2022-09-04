package app

import (
	"context"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	errors "github.com/pkg/errors"
)

const (
	dynamodbTableName = "ingeolmi-dynamodb"
)

func (f *DynamodbTableFields) marshalDynamodbAttributeValue() (map[string]types.AttributeValue, error) {
	var err error

	encodedItem := make(map[string]types.AttributeValue)

	v := reflect.ValueOf(*f)

	for i := 0; i < v.NumField(); i++ {

		fieldName := v.Type().Field(i).Name

		switch v.Field(i).Interface().(type) {

		case int:
			if v.Field(i).Int() == 0 {
				continue
			}

			encodedItem[fieldName], err = attributevalue.Marshal(v.Field(i).Int())
			if err != nil {
				return nil, errors.Wrap(err, "error to masrshal dynamodb attributer value to int")
			}

		case string:
			if v.Field(i).String() == "" {
				continue
			}

			encodedItem[fieldName], err = attributevalue.Marshal(v.Field(i).String())
			if err != nil {
				return nil, errors.Wrap(err, "error to masrshal dynamodb attributer value to string")
			}
		}
	}

	return encodedItem, nil
}

func (i *Injeolmi) PutItemToAWSDynamodb(item map[string]types.AttributeValue) error {
	cli := dynamodb.NewFromConfig(*i.awsClientConfig)

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

func (i *Injeolmi) GetItemToAWSDynamodb(item map[string]types.AttributeValue) (*map[string]types.AttributeValue, error) {
	cli := dynamodb.NewFromConfig(*i.awsClientConfig)

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
