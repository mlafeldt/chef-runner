// Package exec runs external commands. It's a wrapper around Go's os/exec
// package that allows to stub command execution for testing.
package exec

import (
	"os"
	goexec "os/exec"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
)

// The RunnerFunc type is an adapter to use any function for running commands.
type RunnerFunc func(args []string) error

// DefaultRunner is the default function used to run commands. It calls
// os/exec.Run so that stdout and stderr are written to the terminal.
// DefaultRunner also logs all executed commands.
func DefaultRunner(args []string) error {
	log.Debugf("exec: %s\n", strings.Join(args, " "))
	cmd := goexec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var runnerFunc = DefaultRunner

// SetRunnerFunc registers the function f to run all future commands.
func SetRunnerFunc(f RunnerFunc) {
	runnerFunc = f
}

// RunCommand runs the specified command using the currently registered
// RunnerFunc.
func RunCommand(args []string) error {
	return runnerFunc(args)
}
