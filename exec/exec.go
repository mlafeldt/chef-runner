package exec

import (
	"os"
	"os/exec"
)

type CmdRunner func(args []string) error

var cmdRunnerFunc CmdRunner = func(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SetCmdRunnerFunc(f CmdRunner) {
	cmdRunnerFunc = f
}

func ExecuteCommand(args []string) error {
	return cmdRunnerFunc(args)
}
