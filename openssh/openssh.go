package openssh

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mlafeldt/chef-runner/exec"
)

type Client struct {
	Host        string
	User        string
	Port        int
	PrivateKeys []string
	Options     map[string]string
}

// NewClient creates a OpenSSH client from the given host string. The host
// string has the format [user@]hostname[:port]
func NewClient(host string) (*Client, error) {
	var user string
	a := strings.Split(host, "@")
	if len(a) > 1 {
		user = a[0]
		host = a[1]
	}

	var port int
	a = strings.Split(host, ":")
	if len(a) > 1 {
		host = a[0]
		var err error
		if port, err = strconv.Atoi(a[1]); err != nil {
			return nil, errors.New("invalid SSH port")
		}
	}

	c := Client{
		Host: host,
		User: user,
		Port: port,
	}
	return &c, nil
}

func (c Client) String() string {
	return fmt.Sprintf("OpenSSH (host: %s)", c.Host)
}

func (c Client) Command(command string) ([]string, error) {
	if command == "" {
		return nil, errors.New("no command given")
	}
	if c.Host == "" {
		return nil, errors.New("no host given")
	}

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
	return cmd, nil
}

func (c Client) RunCommand(command string) error {
	cmd, err := c.Command(command)
	if err != nil {
		return err
	}
	return exec.RunCommand(cmd)
}
