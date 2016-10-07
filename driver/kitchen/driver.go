// Package kitchen implements a driver based on Test Kitchen.
package kitchen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/mlafeldt/chef-runner/rsync"
)

// Driver is a driver based on Test Kitchen.
type Driver struct {
	Instance    string
	SSHClient   *openssh.Client
	RsyncClient *rsync.Client
}

// This is what `vagrant ssh` uses
var defaultSSHOptions = [...]string{
	"UserKnownHostsFile /dev/null",
	"StrictHostKeyChecking no",
	"PasswordAuthentication no",
	"IdentitiesOnly yes",
	"LogLevel FATAL",
}

type kitchenConfigEntry struct {
	Name string `yaml:"name"`
}

type kitchenConfig struct {
	Platforms []kitchenConfigEntry `yaml:"platforms"`
	Suites    []kitchenConfigEntry `yaml:"suites"`
}

func readInstanceNames() ([]string, error) {
	data, err := ioutil.ReadFile(".kitchen.yml")
	if err != nil {
		return nil, err
	}

	var config kitchenConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var names []string
	for _, suite := range config.Suites {
		for _, platform := range config.Platforms {
			s := suite.Name + "-" + platform.Name
			s = strings.Replace(s, "_", "-", -1)
			s = strings.Replace(s, ".", "", -1)
			names = append(names, s)
		}
	}
	return names, nil
}

func findInstanceName(instance string) (string, error) {
	instanceNames, err := readInstanceNames()
	if err != nil {
		return "", err
	}
	log.Debugf("Kitchen instances = %s\n", instanceNames)

	for _, name := range instanceNames {
		// Return first match
		if strings.Contains(name, instance) {
			return name, nil
		}
	}
	return "", errors.New("Kitchen instance not found")
}

type instanceConfig struct {
	Hostname string
	Username string
	SSHKey   string
	Port     int
}

func (c *instanceConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Hostname string `yaml:"hostname"`
		Username string `yaml:"username"`
		SSHKey   string `yaml:"ssh_key"`
		Port     string `yaml:"port"`
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}
	if aux.Hostname == "" {
		return errors.New("Kitchen config: invalid `hostname`")
	}
	if aux.Username == "" {
		return errors.New("Kitchen config: invalid `username`")
	}
	if aux.SSHKey == "" {
		return errors.New("Kitchen config: invalid `ssh_key`")
	}
	// Test Kitchen stores the port as an string
	port, err := strconv.Atoi(aux.Port)
	if err != nil {
		return errors.New("Kitchen config: invalid `port`")
	}

	c.Hostname = aux.Hostname
	c.Username = aux.Username
	c.SSHKey = aux.SSHKey
	c.Port = port
	return nil
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

	return &config, nil
}

// NewDriver creates a new Test Kitchen driver that communicates with the given
// Test Kitchen instance. Under the hood the instance's YAML configuration is
// parsed to get a working SSH configuration.
func NewDriver(instance string, sshOptions, rsyncOptions []string) (*Driver, error) {
	instance, err := findInstanceName(instance)
	if err != nil {
		return nil, err
	}

	config, err := readInstanceConfig(instance)
	if err != nil {
		return nil, err
	}

	sshOpts := make([]string, len(defaultSSHOptions))
	copy(sshOpts, defaultSSHOptions[:])
	for _, o := range sshOptions {
		sshOpts = append(sshOpts, o)
	}

	sshClient := &openssh.Client{
		Host:        config.Hostname,
		User:        config.Username,
		Port:        config.Port,
		PrivateKeys: []string{config.SSHKey},
		Options:     sshOpts,
	}

	rsyncClient := *rsync.MirrorClient
	rsyncClient.RemoteHost = config.Hostname
	rsyncClient.RemoteShell = sshClient.Shell()
	rsyncClient.Options = rsyncOptions

	return &Driver{instance, sshClient, &rsyncClient}, nil
}

// RunCommand runs the specified command on the Test Kitchen instance.
func (drv Driver) RunCommand(args []string) error {
	return drv.SSHClient.RunCommand(args)
}

// Upload copies files to the Test Kitchen instance.
func (drv Driver) Upload(dst string, src ...string) error {
	return drv.RsyncClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
	return fmt.Sprintf("Test Kitchen driver (instance: %s)", drv.Instance)
}
