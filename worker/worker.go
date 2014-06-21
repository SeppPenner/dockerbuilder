// Package worker implements the worker executing the container build
package worker

import (
	"fmt"
	"github.com/brocaar/dockerbuilder/helpers"
	"github.com/brocaar/dockerbuilder/repository"
	"github.com/brocaar/dockerbuilder/workspace"
	"log"
	"os"
	"os/exec"
	"path"
)

type WorkerTask struct {
	Revision   string
	Repository *repository.Repository
}

// Worker executes the WorkerTask items in the queue.
// This will clone the repository, build the container and on a successful
// build, it will push it to the Docker index.
func Worker(taskQueue chan *WorkerTask) {
	for {
		workerTask := <-taskQueue
		err := cloneOrFetchRepo(workerTask.Repository)
		if err != nil {
			log.Printf("something went wrong while getting or fetching the repository: %s\n", err)
			continue
		}
	}
}

func cloneOrFetchRepo(repo *repository.Repository) error {
	repoClonePath := workspace.GetClonePath(repo)

	if repo.SCM == repository.ScmGit {
		if helpers.PathExists(repoClonePath) {
			// we need to make sure the repository is up-to-date
			log.Printf("fetching repository in: %s\n", repoClonePath)
			cmd := exec.Command("git", "fetch", "--all")
			cmd.Dir = repoClonePath
			return cmd.Run()
		} else {
			// we have to checkout the directory
			repoClonePath = path.Join(repoClonePath, "..")
			err := os.MkdirAll(repoClonePath, 0700)
			if err != nil {
				return err
			}
			log.Printf("cloning repository: %s\n", repo.Url)
			cmd := exec.Command("git", "clone", repo.Url)
			cmd.Dir = repoClonePath
			return cmd.Run()
		}
	}

	return fmt.Errorf("SCM %s is not supported", repo.SCM)
}
