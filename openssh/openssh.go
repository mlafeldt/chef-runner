package openssh

import (
	"github.com/mlafeldt/chef-runner.go/exec"
)

type SSHClient struct {
	Host string
}

func NewSSHClient(host string) *SSHClient {
	return &SSHClient{Host: host}
}

func (c *SSHClient) RunCommand(cmd string) error {
	return exec.RunCommand([]string{"ssh", c.Host, "-c", cmd})
}
