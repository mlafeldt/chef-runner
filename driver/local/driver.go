// Package ssh implements a driver based on OpenSSH.
package local

import (
  "fmt"

  "github.com/mlafeldt/chef-runner/log"
  "github.com/mlafeldt/chef-runner/exec"
  "github.com/mlafeldt/chef-runner/rsync"
)

type Driver struct {
  rsyncClient *rsync.Client
}

// NewDriver creates a new local driver that communicates with the given host.
func NewDriver() (*Driver, error) {
  rsyncClient := rsync.MirrorClient
  return &Driver{rsyncClient}, nil
}

// RunCommand runs the specified command
func (drv Driver) RunCommand(args []string) error {
  cmd := []string{"sh", "-c"}
  if len(args) > 0 {
    cmd = append(cmd, args...)
  }
  log.Infof("%s",args)
  return exec.RunCommand(cmd)
}

func (drv Driver) Upload(dst string, src ...string) error {
  return drv.rsyncClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
  return fmt.Sprintf("local driver")
}
