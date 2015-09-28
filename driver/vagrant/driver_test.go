package vagrant_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/mlafeldt/chef-runner/driver"
	. "github.com/mlafeldt/chef-runner/driver/vagrant"
	"github.com/mlafeldt/chef-runner/util"
)

func TestDriverInterface(t *testing.T) {
	assert.Implements(t, (*driver.Driver)(nil), new(Driver))
}

func TestNewDriver(t *testing.T) {
	util.InDir("../../testdata", func() {
		oldPath := os.Getenv("PATH")
		os.Setenv("PATH", strings.Join([]string{"bin", oldPath},
			string(os.PathListSeparator)))
		defer os.Setenv("PATH", oldPath)

		sshOpts := []string{"LogLevel=debug"}
		rsyncOpts := []string{"--verbose"}
		drv, err := NewDriver("some-machine", sshOpts, rsyncOpts)
		if assert.NoError(t, err) {
			defer os.RemoveAll(".chef-runner")
			assert.Equal(t, "default", drv.SSHClient.Host)
			assert.Equal(t, ".chef-runner/vagrant/machines/some-machine/ssh_config",
				drv.SSHClient.ConfigFile)
			assert.Equal(t, sshOpts, drv.SSHClient.Options)

			assert.Equal(t, "default", drv.RsyncClient.RemoteHost)
			assert.Equal(t, rsyncOpts, drv.RsyncClient.Options)
		}
	})
}

func TestString(t *testing.T) {
	expect := "Vagrant driver (machine: some-machine)"
	actual := Driver{Machine: "some-machine"}.String()
	assert.Equal(t, expect, actual)
}
