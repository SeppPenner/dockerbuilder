// Package helpers provides all kind of "helping" functions.
package helpers

import (
	"os"
)

// PathExists checks if the given path exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
