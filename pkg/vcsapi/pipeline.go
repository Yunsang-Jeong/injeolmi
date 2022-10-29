package vcsapi

import (
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func RunGitlabBranchPipeline(cli gitlab.Client, project_id int, ref string, variables *[]*gitlab.PipelineVariable) (*gitlab.Pipeline, error) {
	opt := &gitlab.CreatePipelineOptions{
		Ref:       gitlab.String(ref),
		Variables: variables,
	}

	// pipeline: https://pkg.go.dev/github.com/xanzy/go-gitlab#Pipeline
	pipeline, _, err := cli.Pipelines.CreatePipeline(project_id, opt)
	if err != nil {
		return nil, errors.Wrap(err, "fail to run branch pipeline")
	}

	return pipeline, nil
}
