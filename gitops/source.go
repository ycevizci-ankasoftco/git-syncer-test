package gitops

import (
	"fmt"
	"os"

	"tr/com/emlakkatilim/git-syncer/config"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneOrPullSourceRepo() (*git.Repository, error) {
	if _, err := os.Stat(config.SourcePath); os.IsNotExist(err) {
		fmt.Println("Cloning source repo...")
		return git.PlainClone(config.SourcePath, false, &git.CloneOptions{
			URL:           config.SourceRepoURL,
			ReferenceName: plumbing.NewBranchReferenceName(config.SourceBranchName),
			SingleBranch:  true,
			Depth:         1,
			Auth:          config.SourceAuth,
		})
	} else {
		fmt.Println("Pulling source repo...")
		repo, err := git.PlainOpen(config.SourcePath)
		if err != nil {
			return nil, err
		}
		w, err := repo.Worktree()
		if err != nil {
			return nil, err
		}
		err = w.Pull(&git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: plumbing.NewBranchReferenceName(config.SourceBranchName),
			Auth:          config.SourceAuth,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}
		return repo, nil
	}
}

func GetLatestCommitHash(repo *git.Repository) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}
