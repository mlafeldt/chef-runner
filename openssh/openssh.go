package openssh

import (
	"github.com/mlafeldt/chef-runner.go/exec"
)

type Client struct {
	Host string
}

func NewClient(host string) *Client {
	return &Client{Host: host}
}

func (c *Client) Command(command string) []string {
	return []string{"ssh", c.Host, command}
}

func (c *Client) RunCommand(command string) error {
	return exec.RunCommand(c.Command(command))
}
