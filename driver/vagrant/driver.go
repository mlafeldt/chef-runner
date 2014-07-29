package vagrant

import (
	"fmt"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/openssh"
)

const (
	DefaultMachine = "default"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

type Driver struct {
	Machine   string
	sshClient *openssh.Client
}

func NewDriver(machine string) *Driver {
	if machine == "" {
		machine = DefaultMachine
	}
	return &Driver{Machine: machine}
}

func (d *Driver) String() string {
	return fmt.Sprintf("Vagrant driver (machine: %s)", d.Machine)
}

func (d *Driver) SSHConfig() (string, error) {
	config, err := goexec.Command("vagrant", "ssh-config", d.Machine).Output()
	if err != nil {
		return "", err
	}
	return string(config), nil
}

func (d *Driver) WriteSSHConfig(filename string) error {
	log.Debug("Writing current SSH config to", filename)
	config, err := d.SSHConfig()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(config), 0644)
}

func (d *Driver) SSHClient() (*openssh.Client, error) {
	if d.sshClient != nil {
		return d.sshClient, nil
	}
	// TODO: reuse existing config file, but make sure it's still valid
	configFile := path.Join(".vagrant", "machines", d.Machine, "ssh_config")
	if err := d.WriteSSHConfig(configFile); err != nil {
		return nil, err
	}
	d.sshClient = &openssh.Client{
		Host:       "default",
		ConfigFile: configFile,
	}
	return d.sshClient, nil
}

func (d *Driver) RunCommand(command string) error {
	sshClient, err := d.SSHClient()
	if err != nil {
		return err
	}
	return sshClient.RunCommand(command)
}
