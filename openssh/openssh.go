package openssh

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/mlafeldt/chef-runner/exec"
)

type Client struct {
	Host        string
	User        string
	Port        int
	PrivateKeys []string
	Options     map[string]string
}

func NewClient(host string) *Client {
	return &Client{Host: host}
}

func (c Client) String() string {
	return fmt.Sprintf("OpenSSH (host: %s)", c.Host)
}

func (c Client) SSHCommand(command string) []string {
	cmd := []string{"ssh"}

	if c.User != "" {
		cmd = append(cmd, "-l", c.User)
	}

	if c.Port != 0 {
		cmd = append(cmd, "-p", strconv.Itoa(c.Port))
	}

	for _, pk := range c.PrivateKeys {
		cmd = append(cmd, "-i", pk)
	}

	// Sort options by name before using them
	var optionNames []string
	for k := range c.Options {
		optionNames = append(optionNames, k)
	}
	sort.Strings(optionNames)
	for _, k := range optionNames {
		cmd = append(cmd, "-o", k+"="+c.Options[k])
	}

	cmd = append(cmd, c.Host, command)
	return cmd
}

func (c Client) RunCommand(command string) error {
	return exec.RunCommand(c.SSHCommand(command))
}
