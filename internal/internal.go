package internal

const (
	DynamodbTablePK = "PipelineID"
)

type DynamodbTableSchema struct {
	PipelineID       int    `dynamodbav:"PipelineID"`
	CommentID        int    `dynamodbav:"CommentID"`
	MergeRequestsIID int    `dynamodbav:"MergeRequestsIID"`
	ActionType       string `dynamodbav:"ActionType"`
	ActionOptions    string `dynamodbav:"ActionOptions"`
	CommentString    string `dynamodbav:"CommentString"`
}
