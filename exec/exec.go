package exec

import (
	"os"
	goexec "os/exec"
)

type RunnerFunc func(args []string) error

func DefaultRunner(args []string) error {
	cmd := goexec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var runnerFunc = DefaultRunner

func SetRunnerFunc(f RunnerFunc) {
	runnerFunc = f
}

func RunCommand(args []string) error {
	return runnerFunc(args)
}
