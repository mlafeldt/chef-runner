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

func (c *SSHClient) RunCommand(command string) error {
	cmd := []string{"ssh", c.Host, command}
	return exec.RunCommand(cmd)
}
