package vagrant

import (
	"os"

	"github.com/mlafeldt/chef-runner.go/exec"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

func RunCommand(machine, cmd string) error {
	if machine == "" {
		machine = "default"
	}
	return exec.RunCommand([]string{"vagrant", "ssh", machine, "-c", cmd})
}

type SSHClient struct {
	Machine string
}

func NewSSHClient(machine string) *SSHClient {
	return &SSHClient{Machine: machine}
}

func (c *SSHClient) RunCommand(cmd string) error {
	return RunCommand(c.Machine, cmd)
}
