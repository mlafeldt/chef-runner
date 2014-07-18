package openssh_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expect := "OpenSSH (host: some-host)"
	actual := openssh.NewClient("some-host").String()
	assert.Equal(t, expect, actual)
}

func TestSSHCommand(t *testing.T) {
	expect := []string{"ssh", "some-host", "uname -a"}
	actual := openssh.NewClient("some-host").SSHCommand("uname -a")
	assert.Equal(t, expect, actual)
}
