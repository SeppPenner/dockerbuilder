// Package workspace contains functions related to preparing the workspace
// for building containers.
package workspace

import (
	"github.com/brocaar/dockerbuilder/config"
	"log"
	"os"
	"os/exec"
)

// Prepare prepares the workspace environment and does some checks.
func Prepare(config *config.Configuration) {
	log.Println("preparing workspace")
	var err error
	var out []byte

	out, err = exec.Command("git", "version").Output()
	if err != nil {
		log.Fatalf("could not execute 'git version': %s", err)
	}
	log.Print(string(out))

	out, err = exec.Command("docker", "version").Output()
	if err != nil {
		log.Fatalf("could not execute 'docker version': %s", err)
	}
	log.Printf("docker version:\n---\n%s---\n", out)

	err = os.Mkdir(config.GetClonePath(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("could not create clone path: %s", config.GetClonePath())
	}

	err = os.Mkdir(config.GetBuildPath(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("cound not create build path: %s", config.GetBuildPath())
	}
}
