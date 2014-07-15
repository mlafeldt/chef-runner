package openssh_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner.go/openssh"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := openssh.NewClient("some-host")
	if assert.NotNil(t, client) {
		assert.Equal(t, "some-host", client.Host)
	}
}

func TestSSHCommand(t *testing.T) {
	expect := []string{"ssh", "some-host", "uname -a"}
	actual := openssh.NewClient("some-host").SSHCommand("uname -a")
	assert.Equal(t, expect, actual)
}
