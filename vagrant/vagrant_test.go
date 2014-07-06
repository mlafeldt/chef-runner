package vagrant_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/vagrant"
	"github.com/stretchr/testify/assert"
)

var lastCmd string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = strings.Join(args, " ")
		return nil
	})
}

func TestRunCommand_DefaultMachine(t *testing.T) {
	err := vagrant.RunCommand("", "uname -a")
	if assert.NoError(t, err) {
		assert.Equal(t, "vagrant ssh default -c uname -a", lastCmd)
	}
}

func TestRunCommand_CustomMachine(t *testing.T) {
	err := vagrant.RunCommand("some-machine", "uname -a")
	if assert.NoError(t, err) {
		assert.Equal(t, "vagrant ssh some-machine -c uname -a", lastCmd)
	}
}
