package workspace

import (
	c "github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/repository"
	"testing"
)

func init() {
	config, _ := c.GetConfiguration()
	SetConfig(config)
}

// TestGetCloneBasePath tests the GetCloneBasePath function.
func TestGetCloneBasePath(t *testing.T) {
	cloneBasePath := GetCloneBasePath()
	if cloneBasePath != "/tmp/clones" {
		t.Errorf("expected: /tmp/clones, got: %s", cloneBasePath)
	}
}

// TestGetBuildBasePath tests the GetBuildBasePath function.
func TestGetBuildBasePath(t *testing.T) {
	buildBasePath := GetBuildBasePath()
	if buildBasePath != "/tmp/builds" {
		t.Errorf("expected: /tmp/builds, got: %s", buildBasePath)
	}
}

// TestGetClonePath tests the GetClonePath function.
func TestGetClonePath(t *testing.T) {
	repo := repository.NewRepository(repository.HostGitHub, "brocaar", "dockerbuilder", repository.ScmGit)
	clonePath := GetClonePath(repo)
	expected := "/tmp/clones/github.com/brocaar/dockerbuilder"
	if clonePath != expected {
		t.Errorf("expected: %s, got: %s", expected, clonePath)
	}
}

// TestGetBuildPath tests the GetBuildPath function.
func TestGetBuildPath(t *testing.T) {
	repo := repository.NewRepository(repository.HostGitHub, "brocaar", "dockerbuilder", repository.ScmGit)
	buildPath := GetBuildPath(repo)
	expected := "/tmp/builds/github.com/brocaar/dockerbuilder"
	if buildPath != expected {
		t.Errorf("expected: %s, got: %s", expected, buildPath)
	}
}
