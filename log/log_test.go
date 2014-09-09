package log_test

import (
	"os"

	. "github.com/mlafeldt/chef-runner/log"
)

func init() {
	DisableColor()
}

func ExampleDebug() {
	Debug("some debug message")
	// Output:
	// DEBUG: some debug message
}

func ExampleDebugf() {
	s := "debug"
	Debugf("some %s message", s)
	// Output:
	// DEBUG: some debug message
}

func ExampleInfo() {
	Info("some info message")
	// Output:
	// INFO: some info message
}

func ExampleInfof() {
	s := "info"
	Infof("some %s message", s)
	// Output:
	// INFO: some info message
}

func ExampleWarn() {
	Warn("some warning message")
	// Output:
	// WARNING: some warning message
}

func ExampleWarnf() {
	s := "warning"
	Warnf("some %s message", s)
	// Output:
	// WARNING: some warning message
}

func ExampleError() {
	os.Stderr = os.Stdout
	Error("some error message")
	// Output:
	// ERROR: some error message
}

func ExampleErrorf() {
	os.Stderr = os.Stdout
	s := "error"
	Errorf("some %s message", s)
	// Output:
	// ERROR: some error message
}

func ExampleSetLevel() {
	defer SetLevel(LevelDebug)
	SetLevel(LevelInfo)

	Debug("some debug message")
	Info("some info message")
	Warn("some warning message")
	// Output:
	// INFO: some info message
	// WARNING: some warning message
}
