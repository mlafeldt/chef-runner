package exec_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	. "github.com/mlafeldt/chef-runner/exec"
)

func TestRunCommand_Success(t *testing.T) {
	err := RunCommand([]string{"go", "version"})
	assert.NoError(t, err)
}

func TestRunCommand_Failure(t *testing.T) {
	err := RunCommand([]string{"go", "some-unknown-subcommand"})
	assert.EqualError(t, err, "exit status 2")
}

func TestRunCommand_Func(t *testing.T) {
	defer SetRunnerFunc(DefaultRunner)

	var lastCmd string
	SetRunnerFunc(func(args []string) error {
		lastCmd = strings.Join(args, " ")
		return nil
	})

	err := RunCommand([]string{"some", "test", "command"})
	if assert.NoError(t, err) {
		assert.Equal(t, "some test command", lastCmd)
	}
}
