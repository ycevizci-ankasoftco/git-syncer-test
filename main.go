package main

import (
	"fmt"
	"time"

	"tr/com/emlakkatilim/git-syncer/config"
	"tr/com/emlakkatilim/git-syncer/gitops"
)

func main() {
	var lastHash string

	for {
		sourceRepo, err := gitops.CloneOrPullSourceRepo()
		if err != nil {
			fmt.Println("Error syncing source repo:", err)
			time.Sleep(config.PollInterval)
			continue
		}

		hash, err := gitops.GetLatestCommitHash(sourceRepo)
		if err != nil {
			fmt.Println("Error reading commit hash:", err)
			time.Sleep(config.PollInterval)
			continue
		}

		if hash != lastHash {
			fmt.Println("New commit detected:", hash)

			_, err := gitops.InitOrOpenTargetRepo()
			if err != nil {
				fmt.Println("Error preparing target repo:", err)
				time.Sleep(config.PollInterval)
				continue
			}

			err = gitops.CopyFiles(config.SourcePath, config.TargetPath)
			if err != nil {
				fmt.Println("Error copying files:", err)
				time.Sleep(config.PollInterval)
				continue
			}

			err = gitops.CommitAndPushTargetRepo()
			if err != nil {
				fmt.Println("Error pushing to target repo:", err)
			}

			lastHash = hash
		} else {
			fmt.Println("No new commits.")
		}

		time.Sleep(config.PollInterval)
	}
}
