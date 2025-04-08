package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var (
	repoURL      = os.Getenv("GIT_REPO_URL")
	accessToken  = os.Getenv("GIT_ACCESS_TOKEN")
	branchName   = os.Getenv("GIT_BRANCH")
	localPath    = "/tmp/git-watcher-repo"
	pollInterval = 60 * time.Second
	scriptPath   = "script.sh"
	auth         = &http.BasicAuth{
		Username: "abdurrahman-osman", // Can be anything, often "git" or "token"
		Password: accessToken,         // The actual PAT
	}
)

func cloneOrPullRepo() (*git.Repository, error) {
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		fmt.Println("Cloning repo... %s", repoURL)
		return git.PlainClone(localPath, false, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.NewBranchReferenceName(branchName),
			SingleBranch:  true,
			Depth:         1,
			Auth:          auth,
		})
	} else {
		fmt.Println("Pulling repo...")
		repo, err := git.PlainOpen(localPath)
		if err != nil {
			return nil, err
		}
		w, err := repo.Worktree()
		if err != nil {
			return nil, err
		}
		err = w.Pull(&git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: plumbing.NewBranchReferenceName(branchName),
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

func runScript() error {
	fmt.Println("Running script:", scriptPath)
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var lastHash string

	for {
		repo, err := cloneOrPullRepo()
		if err != nil {
			fmt.Println("Error syncing repo:", err)
			time.Sleep(pollInterval)
			continue
		}

		hash, err := getLatestCommitHash(repo)
		if err != nil {
			fmt.Println("Error reading commit hash:", err)
			time.Sleep(pollInterval)
			continue
		}

		if hash != lastHash {
			fmt.Println("New commit detected:", hash)
			err := runScript()
			if err != nil {
				fmt.Println("Script failed:", err)
			}
			lastHash = hash
		} else {
			fmt.Println("No new commits.")
		}

		time.Sleep(pollInterval)
	}
}
