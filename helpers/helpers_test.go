package helpers

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestPathExists tests the behaviour of TestPathExists.
func TestPathExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("/tmp", "unittest")
	if err != nil {
		t.Errorf("creating temp dir failed: %s", err)
	}
	if PathExists(tempDir) != true {
		t.Errorf("was expecting function to return true (dir: %s)", tempDir)
	}

	err = os.Remove(tempDir)
	if err != nil {
		t.Errorf("removing dir failed: %s", err)
	}
	if PathExists(tempDir) != false {
		t.Errorf("was expecting function to return false (dir: %s)", tempDir)
	}
}
