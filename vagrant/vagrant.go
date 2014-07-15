package vagrant

import (
	"os"

	"github.com/mlafeldt/chef-runner.go/exec"
)

const (
	DefaultMachine = "default"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

type Client struct {
	Machine string
}

var DefaultClient = &Client{Machine: DefaultMachine}

func NewClient(machine string) *Client {
	if machine == "" {
		machine = DefaultMachine
	}
	return &Client{Machine: machine}
}

func (c Client) SSHCommand(command string) []string {
	return []string{"vagrant", "ssh", c.Machine, "-c", command}
}

func (c Client) RunCommand(command string) error {
	return exec.RunCommand(c.SSHCommand(command))
}
