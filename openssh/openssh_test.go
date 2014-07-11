package openssh_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/openssh"
	"github.com/stretchr/testify/assert"
)

var lastCmd string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = strings.Join(args, " ")
		return nil
	})
}

func TestNewSSHClient(t *testing.T) {
	client := openssh.NewSSHClient("some-host")
	if assert.NotNil(t, client) {
		assert.Equal(t, "some-host", client.Host)
	}
}

func TestRunCommand(t *testing.T) {
	err := openssh.NewSSHClient("some-host").RunCommand("uname -a")
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh some-host uname -a", lastCmd)
	}
}
