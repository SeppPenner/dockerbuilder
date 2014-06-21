package workspace

import (
	c "github.com/brocaar/dockerbuilder/config"
	"testing"
)

func init() {
	config, _ := c.GetConfiguration()
	SetConfig(config)
}

// TestGetClonePath tests the GetClonePath function.
func TestGetClonePath(t *testing.T) {
	clonePath := GetClonePath()
	if clonePath != "/tmp/repositories" {
		t.Errorf("expected: /tmp/repositories, got: %s", clonePath)
	}
}

// TestGetBuildPath tests the GetBuildPath function.
func TestGetBuildPath(t *testing.T) {
	buildPath := GetBuildPath()
	if buildPath != "/tmp/builds" {
		t.Errorf("expected: /tmp/builds, got: %s", buildPath)
	}
}
