package vagrant_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner.go/vagrant"
	"github.com/stretchr/testify/assert"
)

func TestDefaultClient(t *testing.T) {
	client := vagrant.DefaultClient
	if assert.NotNil(t, client) {
		assert.Equal(t, "default", client.Machine)
	}
}

func TestNewClient(t *testing.T) {
	client := vagrant.NewClient("some-machine")
	if assert.NotNil(t, client) {
		assert.Equal(t, "some-machine", client.Machine)
	}
}

func TestSSHCommand_DefaultMachine(t *testing.T) {
	expect := []string{"vagrant", "ssh", "default", "-c", "uname -a"}
	actual := vagrant.DefaultClient.SSHCommand("uname -a")
	assert.Equal(t, expect, actual)
}

func TestSSHCommand_CustomMachine(t *testing.T) {
	expect := []string{"vagrant", "ssh", "some-machine", "-c", "uname -a"}
	actual := vagrant.NewClient("some-machine").SSHCommand("uname -a")
	assert.Equal(t, expect, actual)
}
