package kitchen_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/mlafeldt/chef-runner/driver"
	. "github.com/mlafeldt/chef-runner/driver/kitchen"
	"github.com/mlafeldt/chef-runner/util"
)

func TestDriverInterface(t *testing.T) {
	assert.Implements(t, (*driver.Driver)(nil), new(Driver))
}

func TestNewDriver(t *testing.T) {
	util.InDir("../../testdata", func() {
		sshOpts := []string{"LogLevel=debug"}
		rsyncOpts := []string{"--verbose"}
		drv, err := NewDriver("default-ubuntu-1404", sshOpts, rsyncOpts)
		if assert.NoError(t, err) {
			assert.Equal(t, "127.0.0.1", drv.SSHClient.Host)
			assert.Equal(t, 2222, drv.SSHClient.Port)
			assert.Equal(t, "vagrant", drv.SSHClient.User)
			assert.Equal(t, "/Users/mlafeldt/.vagrant.d/insecure_private_key",
				drv.SSHClient.PrivateKeys[0])
			assert.Equal(t, 6, len(drv.SSHClient.Options))
			assert.Equal(t, "LogLevel=debug", drv.SSHClient.Options[5])

			assert.Equal(t, "127.0.0.1", drv.RsyncClient.RemoteHost)
			assert.Equal(t, rsyncOpts, drv.RsyncClient.Options)
		}
	})
}

func TestString(t *testing.T) {
	expect := "Test Kitchen driver (instance: some-instance)"
	actual := Driver{Instance: "some-instance"}.String()
	assert.Equal(t, expect, actual)
}
