// Package rsync provides a wrapper around the fast rsync file copying tool.
package rsync

import (
	"errors"

	"github.com/mlafeldt/chef-runner/exec"
)

// A Client is an rsync client. It allows you to copy files from one location to
// another using rsync and supports the tool's most useful command-line options.
type Client struct {
	// Archive, if true, enables archive mode.
	Archive bool

	// Delete, if true, deletes extraneous files from destination directories.
	Delete bool

	// Compress, if true, compresses file data during the transfer.
	Compress bool

	// Verbose, if true, increases rsync's verbosity.
	Verbose bool

	// Exclude contains files to be excluded from the transfer.
	Exclude []string

	// RemoteShell specifies the remote shell to use, e.g. ssh.
	RemoteShell string

	// RemoteHost specifies the remote host to copy files to/from.
	RemoteHost string
}

// DefaultClient is a usable rsync client without any options enabled.
var DefaultClient = &Client{}

// MirrorClient is an rsync client configured to mirror files and directories.
var MirrorClient = &Client{
	Archive:  true,
	Delete:   true,
	Compress: true,
	Verbose:  true,
}

// Command returns the rsync command that will be executed when Copy is called.
func (c Client) Command(dst string, src ...string) ([]string, error) {
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

	if c.Compress {
		cmd = append(cmd, "--compress")
	}

	if c.Verbose {
		cmd = append(cmd, "--verbose")
	}

	for _, x := range c.Exclude {
		cmd = append(cmd, "--exclude", x)
	}

	// FIXME: Only copies files to a remote host, not the other way around.
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

// Copy uses rsync to copy one or more src files to dst.
func (c Client) Copy(dst string, src ...string) error {
	cmd, err := c.Command(dst, src...)
	if err != nil {
		return err
	}
	return exec.RunCommand(cmd)
}
