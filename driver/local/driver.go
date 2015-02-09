// Package local implements a driver that provisions the host system.
package local

import (
	"fmt"
	"os"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/rsync"
)

// Driver is a driver for the host system.
type Driver struct {
	Hostname string
}

// NewDriver creates a new driver that provisions the host system.
func NewDriver() (*Driver, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return &Driver{hostname}, nil
}

// RunCommand runs the specified command on the host system.
func (drv Driver) RunCommand(args []string) error {
	return exec.RunCommand(args)
}

// Upload copies files to the right place on the host system.
func (drv Driver) Upload(dst string, src ...string) error {
	return rsync.MirrorClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
	return fmt.Sprintf("Local driver (hostname: %s)", drv.Hostname)
}
