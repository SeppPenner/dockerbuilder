/*
Package config holds the configuration variables used in this app.

Default variables can be overridden by setting environment variables.
Envirnment variables are uppercase, prefixed by BUILDER_. For example:
BUILDER_DEBUG=True, BUILDER_WORKDIR=/home/builder, ...

*/
package config

import (
	"github.com/kelseyhightower/envconfig"
	"runtime"
)

type Configuration struct {
	WorkDir       string
	Debug         bool
	NumWorkers    int
	BindAddress   string
	TaskQueueSize int
}

// GetConfiguration creates a new instance of the Configuration struct,
// looks for matching environment variables, and returns.
func GetConfiguration() (*Configuration, error) {
	config := &Configuration{
		Debug:         false,
		WorkDir:       "/tmp",
		NumWorkers:    runtime.NumCPU(),
		BindAddress:   "0.0.0.0:5000",
		TaskQueueSize: 100000,
	}

	err := envconfig.Process("builder", config)
	return config, err
}
