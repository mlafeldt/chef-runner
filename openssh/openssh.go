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

func (c *Client) RunCommand(command string) error {
	cmd := []string{"ssh", c.Host, command}
	return exec.RunCommand(cmd)
}
