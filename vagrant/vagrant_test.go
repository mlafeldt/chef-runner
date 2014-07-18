package vagrant_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/vagrant"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expect := "Vagrant (machine: some-machine)"
	actual := vagrant.NewClient("some-machine").String()
	assert.Equal(t, expect, actual)
}

func TestSSHCommand_DefaultMachine(t *testing.T) {
	expect := []string{"vagrant", "ssh", "default", "-c", "uname -a"}
	actual := vagrant.NewClient("").SSHCommand("uname -a")
	assert.Equal(t, expect, actual)
}

func TestSSHCommand_CustomMachine(t *testing.T) {
	expect := []string{"vagrant", "ssh", "some-machine", "-c", "uname -a"}
	actual := vagrant.NewClient("some-machine").SSHCommand("uname -a")
	assert.Equal(t, expect, actual)
}
