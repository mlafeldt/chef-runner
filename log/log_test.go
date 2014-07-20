package log_test

import (
	"os"

	"github.com/mlafeldt/chef-runner/log"
)

func init() {
	log.DisableColor()
}

func ExampleDebug() {
	log.Debug("some debug message")
	// Output:
	// DEBUG: some debug message
}

func ExampleInfo() {
	log.Info("some info message")
	// Output:
	// INFO: some info message
}

func ExampleWarn() {
	log.Warn("some warning message")
	// Output:
	// WARNING: some warning message
}

func ExampleError() {
	os.Stderr = os.Stdout
	log.Error("some error message")
	// Output:
	// ERROR: some error message
}

func ExampleSetLevel() {
	defer log.SetLevel(log.LevelDebug)
	log.SetLevel(log.LevelInfo)

	log.Debug("some debug message")
	log.Info("some info message")
	log.Warn("some warning message")
	// Output:
	// INFO: some info message
	// WARNING: some warning message
}
