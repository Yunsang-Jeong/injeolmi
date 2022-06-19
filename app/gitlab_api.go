package app

import (
	"os"

	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type client struct {
	token         string
	webhookSecret string
	cli           *gitlab.Client
	response      map[string]interface{}
}

func (c *client) init() error {
	var err error

	c.token = os.Getenv("GITLAB_TOKEN")
	c.webhookSecret = os.Getenv("GITLAB_WEBHOOK_SECRET")

	if c.token == "" {
		return errors.New("can't load gitlab token")
	} else if c.webhookSecret == "" {
		return errors.New("can't load gitlab webhook secret")
	}

	c.cli, err = gitlab.NewClient(c.token)
	if err != nil {
		return errors.Wrap(err, "can't create a new gitlab client")
	}

	c.response = make(map[string]interface{})

	return nil
}

func (c *client) WriteComment(project_id int, merge_request_iid int, contents string) error {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	note, response, err := c.cli.Notes.CreateMergeRequestNote(project_id, merge_request_iid, opt)
	if err != nil {
		return errors.Wrap(err, "fail to write comment")
	}

	c.response["note"] = note
	c.response["response"] = response

	return nil
}

func (c *client) RunBranchPipeline(project_id int, ref string, variables *[]*gitlab.PipelineVariable) error {
	opt := &gitlab.CreatePipelineOptions{
		Ref:       gitlab.String(ref),
		Variables: variables,
	}

	pipeline, response, err := c.cli.Pipelines.CreatePipeline(project_id, opt)
	if err != nil {
		return errors.Wrap(err, "fail to run branch pipeline")
	}

	c.response["pipeline"] = pipeline
	c.response["response"] = response

	return nil
}
