package main

import (
	"github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/workspace"
	"log"
)

func main() {
	config, err := config.GetConfiguration()
	if err != nil {
		log.Fatal(err.Error())
	}

	workspace.Prepare(config)
}
