package exec_test

import (
	"strings"
	"testing"

	. "github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Be quiet during testing
	log.SetLevel(log.LevelWarn)
}

func TestRunCommand_Success(t *testing.T) {
	err := RunCommand([]string{"bash", "-c", "echo foo | grep -q foo"})
	assert.NoError(t, err)
}

func TestRunCommand_Failure(t *testing.T) {
	err := RunCommand([]string{"bash", "-c", "echo foo | grep -q bar"})
	assert.EqualError(t, err, "exit status 1")
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
