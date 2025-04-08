package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var (
	sourceRepoURL = os.Getenv("GIT_REPO_URL")
	targetRepoURL = os.Getenv("TARGET_REPO_URL")
	accessToken   = os.Getenv("GIT_ACCESS_TOKEN")
	branchName    = os.Getenv("GIT_BRANCH")

	sourcePath = "/tmp/git-source"
	targetPath = "/tmp/git-target"

	pollInterval = 60 * time.Second

	auth = &http.BasicAuth{
		Username: "abdurrahman-osman",
		Password: accessToken,
	}

	targetAuth = &http.BasicAuth{
		Username: "ycevizci-ankasoftco",
		Password: accessToken,
	}
)

func cloneOrPull(path, url string, auth *http.BasicAuth, branch string) (*git.Repository, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Cloning repo:", url)
		return git.PlainClone(path, false, &git.CloneOptions{
			URL:           url,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
			Depth:         1,
			Auth:          auth,
		})
	} else {
		fmt.Println("Pulling latest changes for:", url)
		repo, err := git.PlainOpen(path)
		if err != nil {
			return nil, err
		}
		w, err := repo.Worktree()
		if err != nil {
			return nil, err
		}
		err = w.Pull(&git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			Auth:          auth,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}
		return repo, nil
	}
}

func getLatestCommitHash(repo *git.Repository) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}

func copyFiles(source, dest string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(path, destPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}

func initOrOpenTargetRepo() (*git.Repository, error) {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Println("Cloning target repo...")
		return git.PlainClone(targetPath, false, &git.CloneOptions{
			URL:           targetRepoURL,
			ReferenceName: plumbing.NewBranchReferenceName(branchName),
			SingleBranch:  true,
			Auth:          targetAuth,
		})
	}
	return git.PlainOpen(targetPath)
}

func commitAndPushTargetRepo() error {
	repo, err := git.PlainOpen(targetPath)
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
		Auth:       targetAuth,
		RemoteName: "origin",
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Println("Successfully pushed to target repo.")
	return nil
}

func main() {
	var lastHash string

	for {
		sourceRepo, err := cloneOrPull(sourcePath, sourceRepoURL, auth, branchName)
		if err != nil {
			fmt.Println("Error cloning source repo:", err)
			time.Sleep(pollInterval)
			continue
		}

		hash, err := getLatestCommitHash(sourceRepo)
		if err != nil {
			fmt.Println("Error getting commit hash:", err)
			time.Sleep(pollInterval)
			continue
		}

		if hash != lastHash {
			fmt.Println("New commit detected:", hash)

			_, err := initOrOpenTargetRepo()
			if err != nil {
				fmt.Println("Error preparing target repo:", err)
				time.Sleep(pollInterval)
				continue
			}

			err = copyFiles(sourcePath, targetPath)
			if err != nil {
				fmt.Println("Error copying files:", err)
				time.Sleep(pollInterval)
				continue
			}

			err = commitAndPushTargetRepo()
			if err != nil {
				fmt.Println("Error pushing to target repo:", err)
			}

			lastHash = hash
		} else {
			fmt.Println("No new commits.")
		}

		time.Sleep(pollInterval)
	}
}
