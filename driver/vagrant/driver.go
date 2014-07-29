// Package vagrant implements the driver.Driver interface. Under the hood
// `vagrant ssh-config` is executed to get a working SSH configuration for the
// current Vagrant machine.
package vagrant

import (
	"fmt"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/mlafeldt/chef-runner/rsync"
)

const (
	DefaultMachine = "default"
)

type Driver struct {
	machine     string
	sshClient   *openssh.Client
	rsyncClient *rsync.Client
}

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

func NewDriver(machine string) (*Driver, error) {
	if machine == "" {
		machine = DefaultMachine
	}

	// TODO: reuse existing config file, but make sure it's still valid
	log.Debug("Asking Vagrant for SSH config")
	config, err := goexec.Command("vagrant", "ssh-config", machine).Output()
	if err != nil {
		return nil, err
	}

	configFile := path.Join(".vagrant", "machines", machine, "ssh_config")

	log.Debug("Writing current SSH config to", configFile)
	if err := ioutil.WriteFile(configFile, config, 0644); err != nil {
		return nil, err
	}

	sshClient := &openssh.Client{
		Host:       "default",
		ConfigFile: configFile,
	}

	sshCmd := sshClient.Command("")
	remoteShell := strings.Join(sshCmd[:len(sshCmd)-1], " ")
	rsyncClient := &rsync.Client{
		Archive:     true,
		Delete:      true,
		Verbose:     true,
		RemoteHost:  "default",
		RemoteShell: remoteShell,
	}

	return &Driver{machine, sshClient, rsyncClient}, nil
}

func (drv Driver) String() string {
	return fmt.Sprintf("Vagrant driver (machine: %s)", drv.machine)
}

func (drv Driver) RunCommand(command string) error {
	return drv.sshClient.RunCommand(command)
}

func (drv Driver) Upload(dst string, src ...string) error {
	return drv.rsyncClient.Copy(dst, src...)
}
