package app

import (
	"os"

	errors "github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func (i *Injeolmi) setGitlabClient() error {
	if client, err := gitlab.NewClient(i.gitlabToken); err != nil {
		return errors.Wrap(err, "can't create a new gitlab client")
	} else {
		i.gitlabClient = client
	}

	return nil
}

func (i *Injeolmi) getGitlabToken() error {
	i.gitlabToken = os.Getenv(gitlabTokenEnvKey)
	if i.gitlabToken == "" {
		return errors.New("can't load gitlab token from env")
	}
	return nil
}
