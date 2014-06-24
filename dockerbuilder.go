package main

import (
	"github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/handler"
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
	var taskQueue = make(worker.TaskQueue, config.TaskQueueSize)

	// start worker processes
	log.Printf("starting %d workers\n", config.NumWorkers)
	for i := 0; i < config.NumWorkers; i++ {
		go worker.Worker(taskQueue)
	}

	// add http handler
	gitHubHandler := handler.NewGitHubHandler(taskQueue)
	http.HandleFunc("/github.com/hook", gitHubHandler.Hook)

	log.Printf("starting webserver, listening on: %s", config.BindAddress)
	http.ListenAndServe(config.BindAddress, nil)
}
