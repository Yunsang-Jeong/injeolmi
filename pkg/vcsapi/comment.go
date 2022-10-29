package vcsapi

import (
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func WriteGitlabComment(cli gitlab.Client, project_id int, merge_request_iid int, contents string) (*gitlab.Note, error) {
	opt := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	// note: https://pkg.go.dev/github.com/xanzy/go-gitlab#Note
	note, _, err := cli.Notes.CreateMergeRequestNote(project_id, merge_request_iid, opt)
	if err != nil {
		return nil, errors.Wrap(err, "fail to write comment")
	}

	return note, nil
}

func UpdateGitlabComment(cli gitlab.Client, project_id int, merge_request_iid int, note_id int, contents string) (*gitlab.Note, error) {
	opt := &gitlab.UpdateMergeRequestNoteOptions{
		Body: gitlab.String(contents),
	}

	// note: https://pkg.go.dev/github.com/xanzy/go-gitlab#Note
	note, _, err := cli.Notes.UpdateMergeRequestNote(project_id, merge_request_iid, note_id, opt)
	if err != nil {
		return nil, errors.Wrap(err, "fail to update comment")
	}

	return note, nil
}
