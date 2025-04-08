package config

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	RepoURL       = os.Getenv("GIT_REPO_URL")
	TargetRepoURL = os.Getenv("TARGET_REPO_URL")
	AccessToken   = os.Getenv("GIT_ACCESS_TOKEN")
	BranchName    = os.Getenv("GIT_BRANCH")

	TargetUsername = os.Getenv("TARGET_USERNAME")
	SourceUsername = os.Getenv("SOURCE_USERNAME")
	SourcePath     = "/tmp/git-source"
	TargetPath     = "/tmp/git-target"

	PollInterval = 60 * time.Second

	SourceAuth = &http.BasicAuth{
		Username: SourceUsername,
		Password: AccessToken,
	}

	TargetAuth = &http.BasicAuth{
		Username: TargetUsername,
		Password: AccessToken,
	}
)
