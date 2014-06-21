package main

import (
	"github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/worker"
	"github.com/brocaar/dockerbuilder/workspace"
	"log"
	"net/http"
)

func main() {
	// get configuration from environment variables
	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatal(err.Error())
	}

	// setup and prepare the workspace
	workspace.SetConfig(config)
	workspace.Prepare()

	// make the task queue
	var taskQueue = make(chan *worker.WorkerTask, config.NumWorkers)

	// start worker processes
	log.Printf("starting %d workers\n", config.NumWorkers)
	for i := 0; i < config.NumWorkers; i++ {
		go worker.Worker(taskQueue)
	}

	http.ListenAndServe(config.BindAddress, nil)
}
