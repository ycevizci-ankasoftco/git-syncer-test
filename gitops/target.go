package gitops

import (
	"fmt"
	"time"

	"tr/com/emlakkatilim/git-syncer/config"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func InitOrOpenTargetRepo() (*git.Repository, error) {
	if _, err := git.PlainOpen(config.TargetPath); err == nil {
		return git.PlainOpen(config.TargetPath)
	}
	fmt.Println("Cloning target repo...")
	return git.PlainClone(config.TargetPath, false, &git.CloneOptions{
		URL:           config.TargetRepoURL,
		ReferenceName: plumbing.NewBranchReferenceName(config.TargetBranchName),
		SingleBranch:  true,
		Auth:          config.TargetAuth,
	})
}

func CommitAndPushTargetRepo() error {
	repo, err := git.PlainOpen(config.TargetPath)
	if err != nil {
		return fmt.Errorf("failed to open target repo: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	if status.IsClean() {
		fmt.Println("No changes to commit.")
		return nil
	}

	_, err = w.Commit("Synced changes from source repo", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Git Sync Bot",
			Email: "bot@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	err = repo.Push(&git.PushOptions{
		Auth:       config.TargetAuth,
		RemoteName: "origin",
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Println("Successfully pushed to target repo.")
	return nil
}
