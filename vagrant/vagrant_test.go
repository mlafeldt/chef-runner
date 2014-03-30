package vagrant_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/vagrant"
	"github.com/stretchr/testify/assert"
)

var history []string

func clearHistory() { history = []string{} }

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		history = append(history, strings.Join(args, " "))
		return nil
	})
}

func TestRunCommand_DefaultMachine(t *testing.T) {
	defer clearHistory()
	err := vagrant.RunCommand("", "uname -a")
	if assert.NoError(t, err) && assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "vagrant ssh default -c uname -a", history[0])
	}
}

func TestRunCommand_CustomMachine(t *testing.T) {
	defer clearHistory()
	err := vagrant.RunCommand("some-machine", "uname -a")
	if assert.NoError(t, err) && assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "vagrant ssh some-machine -c uname -a", history[0])
	}
}
