package openssh

import (
	"fmt"

	"github.com/mlafeldt/chef-runner/exec"
)

type Client struct {
	Host string
}

func NewClient(host string) *Client {
	return &Client{Host: host}
}

func (c Client) String() string {
	return fmt.Sprintf("OpenSSH (host: %s)", c.Host)
}

func (c Client) SSHCommand(command string) []string {
	return []string{"ssh", c.Host, command}
}

func (c Client) RunCommand(command string) error {
	return exec.RunCommand(c.SSHCommand(command))
}
