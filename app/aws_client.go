package app

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	errors "github.com/pkg/errors"
)

const (
	awsClientRegion = "ap-northeast-2"
)

type dynamodb_table_attributes struct {
	CommentID     int
	MRID          int
	MRIID         int
	PipelineID    int
	ActionType    string
	ActionOptions string
}

func (i *Injeolmi) setAWSClient() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsClientRegion))
	if err != nil {
		return errors.Wrap(err, "can't create a new aws client")
	}
	i.awsClientConfig = &cfg

	return nil
}
