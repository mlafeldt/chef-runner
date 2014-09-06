// Package ssh implements a driver based on OpenSSH.
package ssh

import (
	"fmt"

	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/mlafeldt/chef-runner/rsync"
)

// Driver is a driver based on SSH.
type Driver struct {
	host        string
	sshClient   *openssh.Client
	rsyncClient *rsync.Client
}

// NewDriver creates a new SSH driver that communicates with the given host.
func NewDriver(host string) (*Driver, error) {
	sshClient, err := openssh.NewClient(host)
	if err != nil {
		return nil, err
	}

	rsyncClient := *rsync.MirrorClient
	rsyncClient.RemoteHost = sshClient.Host
	rsyncClient.RemoteShell = sshClient.Shell()

	return &Driver{host, sshClient, &rsyncClient}, nil
}

// RunCommand runs the specified command on the host.
func (drv Driver) RunCommand(args []string) error {
	return drv.sshClient.RunCommand(args)
}

// Upload copies files to the host.
func (drv Driver) Upload(dst string, src ...string) error {
	return drv.rsyncClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
	return fmt.Sprintf("SSH driver (host: %s)", drv.sshClient.Host)
}
