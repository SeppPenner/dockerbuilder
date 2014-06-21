// Package workspace contains functions related to preparing the workspace
// for building containers.
package workspace

import (
	c "github.com/brocaar/dockerbuilder/config"
	"log"
	"os"
	"os/exec"
	"path"
)

var config *c.Configuration

// SetConfig sets the configuration for this workspace.
func SetConfig(conf *c.Configuration) {
	config = conf
}

// Prepare prepares the workspace environment and does some checks.
func Prepare() {
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

	err = os.Mkdir(GetClonePath(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("could not create clone path: %s", GetClonePath())
	}

	err = os.Mkdir(GetBuildPath(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("cound not create build path: %s", GetBuildPath())
	}
}

// GetClonePath returns the absolute path for cloning the repositories in.
func GetClonePath() string {
	return path.Join(config.WorkDir, "repositories")
}

// GetBuildPath returns the absolute path for building the containers in.
func GetBuildPath() string {
	return path.Join(config.WorkDir, "builds")
}
