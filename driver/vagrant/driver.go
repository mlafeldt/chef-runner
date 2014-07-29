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

type Driver struct {
	machine   string
	sshClient *openssh.Client
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

	drv := Driver{
		machine: machine,
		sshClient: &openssh.Client{
			Host:       "default",
			ConfigFile: configFile,
		},
	}
	return &drv, nil
}

func (drv Driver) String() string {
	return fmt.Sprintf("Vagrant driver (machine: %s)", drv.machine)
}

func (drv Driver) RunCommand(command string) error {
	return drv.sshClient.RunCommand(command)
}
