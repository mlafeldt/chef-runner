package ssh_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/driver"
	. "github.com/mlafeldt/chef-runner/driver/ssh"
	"github.com/stretchr/testify/assert"
)

func TestDriverInterface(t *testing.T) {
	assert.Implements(t, (*driver.Driver)(nil), new(Driver))
}

func TestNewDriver(t *testing.T) {
	sshOpts := []string{"LogLevel=debug"}
	drv, err := NewDriver("some-user@some-host:1234", sshOpts)
	if assert.NoError(t, err) {
		assert.Equal(t, "some-host", drv.SSHClient.Host)
		assert.Equal(t, 1234, drv.SSHClient.Port)
		assert.Equal(t, "some-user", drv.SSHClient.User)
		assert.Equal(t, sshOpts, drv.SSHClient.Options)
		assert.Equal(t, "some-host", drv.RsyncClient.RemoteHost)
	}
}

func TestString(t *testing.T) {
	drv, _ := NewDriver("some-user@some-host:1234", nil)
	assert.Equal(t, "SSH driver (host: some-host)", drv.String())
}
