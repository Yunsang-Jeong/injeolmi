package internal

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Yunsang-Jeong/ingeolmi/pkg/awsapi"
	"github.com/Yunsang-Jeong/ingeolmi/pkg/vcsapi"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (c *GitlabEventHandlerConfig) HandlePipeline(event interface{}) error {
	e := event.(gitlab.PipelineEvent)

	returnedItem, err := awsapi.GetItemToAWSDynamodb(*c.AWSConfig, c.DynamodbTableName, map[string]types.AttributeValue{
		DynamodbTablePK: &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", e.ObjectAttributes.ID)},
	})
	if err != nil {
		return errors.Wrap(err, "fail to get item to AWS DynamoDB")
	}

	unmarshaledReturnedItem := make(map[string]string)
	for itemKey, itemValue := range *returnedItem {
		var v string
		if err := attributevalue.Unmarshal(itemValue, &v); err != nil {
			return errors.Wrap(err, "can't parse unmarshal")
		}

		unmarshaledReturnedItem[itemKey] = v
	}

	mergeRequestsIID, _ := strconv.Atoi(unmarshaledReturnedItem["MergeRequestsIID"])
	commentID, _ := strconv.Atoi(unmarshaledReturnedItem["CommentID"])
	commentString := unmarshaledReturnedItem["CommentString"]

	sort.Slice(e.Builds, func(i, j int) bool {
		return e.Builds[i].ID < e.Builds[j].ID
	})

	for _, build := range e.Builds {
		prefix_emoji := ":white_circle:"

		switch build.Status {
		case "success":
			prefix_emoji = ":black_circle:"
		case "running":
			prefix_emoji = ":large_blue_circle:"
		case "failed":
			prefix_emoji = ":red_circle:"
		case "manual":
			prefix_emoji = ":red_circle:"
		}

		c := fmt.Sprintf("- %s - %s is %s : %s", prefix_emoji, build.Stage, build.Name, build.Status)
		commentString = fmt.Sprintf("%s\n%s", commentString, c)
	}

	if _, err := vcsapi.UpdateGitlabComment(
		*c.GitlabClient,
		e.Project.ID,
		mergeRequestsIID,
		commentID,
		commentString,
	); err != nil {
		return errors.Wrap(err, "fail to update comment")
	}

	return nil
}
