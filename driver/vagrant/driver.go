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

type Client struct {
	Machine   string
	sshClient *openssh.Client
}

func NewClient(machine string) *Client {
	if machine == "" {
		machine = DefaultMachine
	}
	return &Client{Machine: machine}
}

func (c *Client) String() string {
	return fmt.Sprintf("Vagrant (machine: %s)", c.Machine)
}

func (c *Client) SSHConfig() (string, error) {
	config, err := goexec.Command("vagrant", "ssh-config", c.Machine).Output()
	if err != nil {
		return "", err
	}
	return string(config), nil
}

func (c *Client) WriteSSHConfig(filename string) error {
	log.Debug("Writing current SSH config to", filename)
	config, err := c.SSHConfig()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(config), 0644)
}

func (c *Client) SSHClient() (*openssh.Client, error) {
	if c.sshClient != nil {
		return c.sshClient, nil
	}
	// TODO: reuse existing config file, but make sure it's still valid
	configFile := path.Join(".vagrant", "machines", c.Machine, "ssh_config")
	if err := c.WriteSSHConfig(configFile); err != nil {
		return nil, err
	}
	c.sshClient = &openssh.Client{
		Host:       "default",
		ConfigFile: configFile,
	}
	return c.sshClient, nil
}

func (c *Client) RunCommand(command string) error {
	sshClient, err := c.SSHClient()
	if err != nil {
		return err
	}
	return sshClient.RunCommand(command)
}
