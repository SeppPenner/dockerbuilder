// Package worker implements the worker executing the container build
package worker

import (
	"fmt"
	"github.com/brocaar/dockerbuilder/helpers"
	"github.com/brocaar/dockerbuilder/repository"
	"github.com/brocaar/dockerbuilder/workspace"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

// A channel for worker tasks.
type TaskQueue chan *WorkerTask

// A struct containing all the data for a worker task
type WorkerTask struct {
	Revision             string
	DockerIndexNamespace string
	Repository           *repository.Repository
	CleanupContainer     bool
}

// Worker executes the WorkerTask items in the queue.
// This will clone the repository, build the container and on a successful
// build, it will push it to the Docker index.
func Worker(taskQueue TaskQueue) {
	var buildPath string
	var containerName string
	var err error

	for {
		workerTask := <-taskQueue

		func() {
			// get an up-to-date copy of the repository
			err = cloneOrFetchRepo(workerTask.Repository)
			if err != nil {
				log.Printf("something went wrong while getting or fetching the repository: %s\n", err)
				return
			}

			// prepare the build path
			buildPath, err = prepareAndGetBuildPath(workerTask.Repository, workerTask.Revision)
			if err != nil {
				log.Printf("something went wrong while preparing the build path: %s\n", err)
				return
			}

			// cleanup build when the function returns
			defer func() {
				log.Printf("removing build: %s\n", buildPath)
				err = os.RemoveAll(buildPath)
				if err != nil {
					log.Printf("removing build failed: %s\n", err)
				}
			}()

			// build container
			containerName = getContainerName(workerTask.Repository, workerTask.Revision, workerTask.DockerIndexNamespace)
			err = buildContainer(buildPath, containerName)
			if err != nil {
				log.Printf("failed building the container: %s, in :%s, reason: %s", containerName, buildPath, err)
				return
			}

			// cleanup container when the function returns
			if workerTask.CleanupContainer == true {
				defer func() {
					err = removeContainer(containerName)
					if err != nil {
						log.Printf("removing container failed: %s\n", containerName)
					}
				}()
			}

			// push container
			err = pushContainer(containerName)
			if err != nil {
				log.Printf("pushing the container failed: %s\n", err)
				return
			}
		}()
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

func prepareAndGetBuildPath(repo *repository.Repository, revision string) (string, error) {
	var err error
	buildPath := workspace.GetBuildPath(repo)
	clonePath := workspace.GetClonePath(repo)

	// make sure the build path exists
	err = os.MkdirAll(buildPath, 0700)
	if err != nil {
		return "", err
	}

	// create temp directory in the build path (it is possible that there are
	// multiple build for the same repo)
	buildPath, err = ioutil.TempDir(buildPath, "build")
	if err != nil {
		return "", err
	}

	// clone repository and checkout the right revision
	if repo.SCM == repository.ScmGit {
		log.Printf("cloning repo: %s into build dir: %s\n", clonePath, buildPath)
		cmd := exec.Command("git", "clone", clonePath, buildPath)
		err = cmd.Run()
		if err != nil {
			return "", err
		}

		log.Printf("checking out revision: %s, in: %s\n", revision, buildPath)
		cmd = exec.Command("git", "checkout", revision)
		cmd.Dir = buildPath
		err = cmd.Run()
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("SCM %s is not supported", repo.SCM)
	}

	return buildPath, nil
}

func getContainerName(repo *repository.Repository, revision, dockerIndexNamespace string) string {
	if dockerIndexNamespace != "" {
		return fmt.Sprintf("%s/%s:%s", dockerIndexNamespace, repo.Name, revision)
	} else {
		return fmt.Sprintf("%s:%s", repo.Name, revision)
	}
}

func buildContainer(buildPath, containerName string) error {
	log.Printf("building container: %s, in: %s\n", containerName, buildPath)
	cmd := exec.Command("docker", "build", "-t", containerName, ".")
	cmd.Dir = buildPath
	return cmd.Run()
}

func removeContainer(containerName string) error {
	log.Printf("removing container: %s\n", containerName)
	cmd := exec.Command("docker", "rmi", containerName)
	return cmd.Run()
}

func pushContainer(containerName string) error {
	log.Printf("pushing container: %s\n", containerName)
	cmd := exec.Command("docker", "push", containerName)
	return cmd.Run()
}
