package app

import (
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (i *Injeolmi) writeComment(project_id int, merge_request_iid int, contents string) error {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	note, response, err := i.gitlabClient.Notes.CreateMergeRequestNote(project_id, merge_request_iid, opt)
	if err != nil {
		return errors.Wrap(err, "fail to write comment")
	}

	i.gitlabClientResponse["note"] = note
	i.gitlabClientResponse["response"] = response

	return nil
}

func (i *Injeolmi) updateComment(project_id int, merge_request_iid int, note_id int, contents string) error {
	opt := &gitlab.UpdateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	note, response, err := i.gitlabClient.Notes.UpdateMergeRequestNote(project_id, merge_request_iid, note_id, opt)
	if err != nil {
		return errors.Wrap(err, "fail to update comment")
	}

	i.gitlabClientResponse["note"] = note
	i.gitlabClientResponse["response"] = response

	return nil
}

func (i *Injeolmi) runBranchPipeline(project_id int, ref string, variables *[]*gitlab.PipelineVariable) error {
	opt := &gitlab.CreatePipelineOptions{
		Ref:       gitlab.String(ref),
		Variables: variables,
	}

	pipeline, response, err := i.gitlabClient.Pipelines.CreatePipeline(project_id, opt)
	if err != nil {
		return errors.Wrap(err, "fail to run branch pipeline")
	}

	// pipeline.ID

	i.gitlabClientResponse["pipeline"] = pipeline
	i.gitlabClientResponse["response"] = response

	return nil
}
