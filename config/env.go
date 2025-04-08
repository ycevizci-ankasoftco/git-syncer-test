package config

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	SourceRepoURL     = os.Getenv("SOURCE_REPO_URL")
	TargetRepoURL     = os.Getenv("TARGET_REPO_URL")
	SourceAccessToken = os.Getenv("SOURCE_ACCESS_TOKEN")
	TargetAccessToken = os.Getenv("TARGET_ACCESS_TOKEN")
	SourceBranchName  = os.Getenv("SOURCE_BRANCH")
	TargetBranchName  = os.Getenv("TARGET_BRANCH")

	TargetUsername = os.Getenv("TARGET_USERNAME")
	SourceUsername = os.Getenv("SOURCE_USERNAME")
	SourcePath     = "/tmp/git-source"
	TargetPath     = "/tmp/git-target"

	PollInterval = 60 * time.Second

	SourceAuth = &http.BasicAuth{
		Username: SourceUsername,
		Password: SourceAccessToken,
	}

	TargetAuth = &http.BasicAuth{
		Username: TargetUsername,
		Password: SourceAccessToken,
	}
)
