package vcsapi

import (
	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func GetGitlabClient(gitlabToken string) (*gitlab.Client, error) {
	if client, err := gitlab.NewClient(gitlabToken); err != nil {
		return nil, errors.Wrap(err, "can't create a new gitlab client")
	} else {
		return client, nil
	}
}
