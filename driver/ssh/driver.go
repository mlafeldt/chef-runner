// Package ssh implements the driver.Driver interface. The implementation is
// just a simple wrapper around the openssh package.
package ssh

import (
	"fmt"

	"github.com/mlafeldt/chef-runner/openssh"
)

type Driver struct {
	host      string
	sshClient *openssh.Client
}

func NewDriver(host string) (*Driver, error) {
	c, err := openssh.NewClient(host)
	if err != nil {
		return nil, err
	}
	drv := Driver{
		host:      host,
		sshClient: c,
	}
	return &drv, nil
}

func (drv Driver) String() string {
	return fmt.Sprintf("SSH driver (host: %s)", drv.sshClient.Host)
}

func (drv Driver) RunCommand(command string) error {
	return drv.sshClient.RunCommand(command)
}
