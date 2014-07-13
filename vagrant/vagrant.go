package vagrant

import (
	"os"

	"github.com/mlafeldt/chef-runner.go/exec"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

type Client struct {
	Machine string
}

func NewClient(machine string) *Client {
	if machine == "" {
		machine = "default"
	}
	return &Client{Machine: machine}
}

func (c *Client) RunCommand(command string) error {
	cmd := []string{"vagrant", "ssh", c.Machine, "-c", command}
	return exec.RunCommand(cmd)
}
