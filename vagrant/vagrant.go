package vagrant

import (
	"os"

	"github.com/mlafeldt/chef-runner.go/exec"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

type SSHClient struct {
	Machine string
}

func NewSSHClient(machine string) *SSHClient {
	if machine == "" {
		machine = "default"
	}
	return &SSHClient{Machine: machine}
}

func (c *SSHClient) RunCommand(command string) error {
	cmd := []string{"vagrant", "ssh", c.Machine, "-c", command}
	return exec.RunCommand(cmd)
}
