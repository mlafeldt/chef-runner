package rsync

import (
	"errors"

	"github.com/mlafeldt/chef-runner/exec"
)

type Client struct {
	Archive bool
	Delete  bool
	Verbose bool
	Exclude []string

	RemoteShell string
	RemoteHost  string
}

var DefaultClient = &Client{}

func (c Client) Command(src []string, dst string) ([]string, error) {
	if len(src) == 0 {
		return nil, errors.New("no source given")
	}
	if dst == "" {
		return nil, errors.New("no destination given")
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

	if c.RemoteShell != "" {
		if c.RemoteHost == "" {
			return nil, errors.New("no remote host given")
		}
		cmd = append(cmd, "--rsh", c.RemoteShell)
		dst = c.RemoteHost + ":" + dst
	}

	cmd = append(cmd, src...)
	cmd = append(cmd, dst)
	return cmd, nil
}

func (c Client) Copy(src []string, dst string) error {
	cmd, err := c.Command(src, dst)
	if err != nil {
		return err
	}
	return exec.RunCommand(cmd)
}
