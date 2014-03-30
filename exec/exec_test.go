package exec_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/stretchr/testify/assert"
)

func TestRunCommand_Success(t *testing.T) {
	err := exec.RunCommand([]string{"bash", "-c", "echo foo | grep -q foo"})
	assert.NoError(t, err)
}

func TestRunCommand_Failure(t *testing.T) {
	err := exec.RunCommand([]string{"bash", "-c", "echo foo | grep -q bar"})
	assert.EqualError(t, err, "exit status 1")
}

func TestRunCommand_Func(t *testing.T) {
	defer exec.SetRunnerFunc(exec.DefaultRunner)

	var history []string
	exec.SetRunnerFunc(func(args []string) error {
		history = append(history, strings.Join(args, " "))
		return nil
	})

	err := exec.RunCommand([]string{"some", "test", "command"})
	if assert.NoError(t, err) && assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "some test command", history[0])
	}
}
