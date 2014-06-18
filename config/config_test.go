package config

import (
	"testing"
)

// TestGetClonePath tests the GetClonePath function.
func TestGetClonePath(t *testing.T) {
	config, _ := GetConfiguration()
	clonePath := config.GetClonePath()

	if clonePath != "/tmp/repositories" {
		t.Errorf("expected: /tmp/repositories, got: %s", clonePath)
	}
}

// TestGetBuildPath tests the GetBuildPath function.
func TestGetBuildPath(t *testing.T) {
	config, _ := GetConfiguration()
	buildPath := config.GetBuildPath()

	if buildPath != "/tmp/builds" {
		t.Errorf("expected: /tmp/builds, got: %s", buildPath)
	}
}
