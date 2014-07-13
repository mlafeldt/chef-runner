package rsync

import (
	"errors"

	"github.com/mlafeldt/chef-runner.go/exec"
)

type Client struct {
	Archive bool
	Delete  bool
	Verbose bool
	Exclude []string
}

var DefaultClient = &Client{}

func (c *Client) Copy(src []string, dst string) error {
	if len(src) == 0 {
		return errors.New("No source given")
	}
	if dst == "" {
		return errors.New("No destination given")
	}

	cmd := []string{"rsync"}
	if c.Archive {
		cmd = append(cmd, "--archive")
	}
	if c.Delete {
		cmd = append(cmd, "--delete")
	}
	if c.Verbose {
		cmd = append(cmd, "--verbose")
	}
	for _, x := range c.Exclude {
		cmd = append(cmd, "--exclude", x)
	}
	cmd = append(cmd, src...)
	cmd = append(cmd, dst)
	return exec.RunCommand(cmd)
}
