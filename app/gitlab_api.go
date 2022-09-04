package app

import (
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (i *Injeolmi) writeComment(project_id int, merge_request_iid int, contents string) (*gitlab.Note, error) {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	// note: https://pkg.go.dev/github.com/xanzy/go-gitlab#Note
	note, _, err := i.gitlabClient.Notes.CreateMergeRequestNote(project_id, merge_request_iid, opt)
	if err != nil {
		return nil, errors.Wrap(err, "fail to write comment")
	}

	return note, nil
}

func (i *Injeolmi) updateComment(project_id int, merge_request_iid int, note_id int, contents string) (*gitlab.Note, error) {
	opt := &gitlab.UpdateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	// note: https://pkg.go.dev/github.com/xanzy/go-gitlab#Note
	note, _, err := i.gitlabClient.Notes.UpdateMergeRequestNote(project_id, merge_request_iid, note_id, opt)
	if err != nil {
		return nil, errors.Wrap(err, "fail to update comment")
	}

	return note, nil
}

func (i *Injeolmi) runBranchPipeline(project_id int, ref string, variables *[]*gitlab.PipelineVariable) (*gitlab.Pipeline, error) {
	opt := &gitlab.CreatePipelineOptions{
		Ref:       gitlab.String(ref),
		Variables: variables,
	}

	// pipeline: https://pkg.go.dev/github.com/xanzy/go-gitlab#Pipeline
	pipeline, _, err := i.gitlabClient.Pipelines.CreatePipeline(project_id, opt)
	if err != nil {
		return nil, errors.Wrap(err, "fail to run branch pipeline")
	}

	return pipeline, nil
}
