// Package kitchen implements a driver based on Test Kitchen.
package kitchen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/mlafeldt/chef-runner/rsync"
	"gopkg.in/yaml.v1"
)

// Driver is a driver based on Test Kitchen.
type Driver struct {
	instance    string
	sshClient   *openssh.Client
	rsyncClient *rsync.Client
}

type instanceConfig struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	SSHKey   string `yaml:"ssh_key"`
	Port     string `yaml:"port"`
}

func readInstanceConfig(instance string) (*instanceConfig, error) {
	configFile := path.Join(".kitchen", instance+".yml")
	log.Debugf("Kitchen config file = %s\n", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config instanceConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	log.Debugf("Kitchen config = %+v\n", config)

	if config.Hostname == "" {
		return nil, errors.New(configFile + ": invalid `hostname`")
	}
	if config.Username == "" {
		return nil, errors.New(configFile + ": invalid `username`")
	}
	if config.SSHKey == "" {
		return nil, errors.New(configFile + ": invalid `ssh_key`")
	}
	if _, err := strconv.Atoi(config.Port); err != nil {
		return nil, errors.New(configFile + ": invalid `port`")
	}

	return &config, nil
}

// NewDriver creates a new Test Kitchen driver that communicates with the given
// Test Kitchen instance. Under the hood the instance's YAML configuration is
// parsed to get a working SSH configuration.
func NewDriver(instance string) (*Driver, error) {
	config, err := readInstanceConfig(instance)
	if err != nil {
		return nil, err
	}

	// Test Kitchen stores the port as an string
	port, _ := strconv.Atoi(config.Port)

	// This is what `vagrant ssh` uses
	sshOpts := map[string]string{
		"UserKnownHostsFile":     "/dev/null",
		"StrictHostKeyChecking":  "no",
		"PasswordAuthentication": "no",
		"IdentitiesOnly":         "yes",
		"LogLevel":               "FATAL",
	}
	sshClient := &openssh.Client{
		Host:        config.Hostname,
		User:        config.Username,
		Port:        port,
		PrivateKeys: []string{config.SSHKey},
		Options:     sshOpts,
	}

	rsyncClient := rsync.MirrorClient
	rsyncClient.RemoteHost = config.Hostname
	rsyncClient.RemoteShell = sshClient.Shell()

	return &Driver{instance, sshClient, rsyncClient}, nil
}

// RunCommand runs the specified command on the Test Kitchen instance.
func (drv Driver) RunCommand(args []string) error {
	return drv.sshClient.RunCommand(args)
}

// Upload copies files to the Test Kitchen instance.
func (drv Driver) Upload(dst string, src ...string) error {
	return drv.rsyncClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
	return fmt.Sprintf("Test Kitchen driver (instance: %s)", drv.instance)
}
