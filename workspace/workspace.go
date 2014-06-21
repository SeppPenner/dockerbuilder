// Package workspace contains functions related to preparing the workspace
// for building containers.
package workspace

import (
	c "github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/repository"
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

	err = os.Mkdir(GetCloneBasePath(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("could not create clone path: %s", GetCloneBasePath())
	}

	err = os.Mkdir(GetBuildBasePath(), 0700)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("cound not create build path: %s", GetBuildBasePath())
	}
}

// GetCloneBasePath returns the absolute path for cloning the repositories in.
func GetCloneBasePath() string {
	return path.Join(config.WorkDir, "clones")
}

// GetBuildBasePath returns the absolute path for building the containers in.
func GetBuildBasePath() string {
	return path.Join(config.WorkDir, "builds")
}

// GetClonePath returns the absolute path for cloning the repository in.
func GetClonePath(repo *repository.Repository) string {
	return path.Join(GetCloneBasePath(), repo.Host, repo.Owner, repo.Name)
}
