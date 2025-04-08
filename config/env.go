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

	SourcePath = "/tmp/git-source"
	TargetPath = "/tmp/git-target"

	PollInterval = 60 * time.Second

	Auth = &http.BasicAuth{
		Username: "abdurrahman-osman",
		Password: AccessToken,
	}

	TargetAuth = &http.BasicAuth{
		Username: "ycevizci-ankasoftco",
		Password: AccessToken,
	}
)
