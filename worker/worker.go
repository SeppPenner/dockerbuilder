// Package worker implements the worker executing the container build
package worker

import (
	"github.com/brocaar/dockerbuilder/repository"
	"log"
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
		log.Printf("Received task: %+v", workerTask)
	}
}
