package openssh_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/stretchr/testify/assert"
)

var newClientTests = []struct {
	host      string
	client    *openssh.Client
	errString string
}{
	{
		host: "some-host",
		client: &openssh.Client{
			Host: "some-host",
		},
	},
	{
		host: "some-user@some-host",
		client: &openssh.Client{
			Host: "some-host",
			User: "some-user",
		},
	},
	{
		host: "some-host:1234",
		client: &openssh.Client{
			Host: "some-host",
			Port: 1234,
		},
	},
	{
		host: "some-user@some-host:1234",
		client: &openssh.Client{
			Host: "some-host",
			User: "some-user",
			Port: 1234,
		},
	},
	// Check for errors
	{
		host:      "some-host:abc",
		client:    nil,
		errString: "invalid SSH port",
	},
}

func TestNewClient(t *testing.T) {
	for _, test := range newClientTests {
		c, err := openssh.NewClient(test.host)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.client, c)
	}
}

func TestString(t *testing.T) {
	expect := "OpenSSH (host: some-host)"
	c, _ := openssh.NewClient("some-host")
	assert.Equal(t, expect, c.String())
}

var commandTests = []struct {
	client    openssh.Client
	cmd       string
	result    []string
	errString string
}{
	{
		client: openssh.Client{
			Host: "some-host",
		},
		cmd:    "uname -a",
		result: []string{"ssh", "some-host", "uname -a"},
	},
	{
		client: openssh.Client{
			Host: "some-host",
			User: "some-user",
		},
		cmd:    "uname -a",
		result: []string{"ssh", "-l", "some-user", "some-host", "uname -a"},
	},
	{
		client: openssh.Client{
			Host: "some-host",
			Port: 1234,
		},
		cmd:    "uname -a",
		result: []string{"ssh", "-p", "1234", "some-host", "uname -a"},
	},
	{
		client: openssh.Client{
			Host:        "some-host",
			PrivateKeys: []string{"some-key", "another-key"},
		},
		cmd: "uname -a",
		result: []string{"ssh", "-i", "some-key", "-i", "another-key",
			"some-host", "uname -a"},
	},
	{
		client: openssh.Client{
			Host: "some-host",
			Options: map[string]string{
				"SomeOption":    "yes",
				"AnotherOption": "no",
			},
		},
		cmd: "uname -a",
		result: []string{"ssh", "-o", "AnotherOption=no",
			"-o", "SomeOption=yes", "some-host", "uname -a"},
	},
	{
		client: openssh.Client{
			Host:        "some-host",
			User:        "some-user",
			Port:        1234,
			PrivateKeys: []string{"some-key"},
			Options:     map[string]string{"SomeOption": "yes"},
		},
		cmd: "uname -a",
		result: []string{"ssh", "-l", "some-user", "-p", "1234",
			"-i", "some-key", "-o", "SomeOption=yes", "some-host",
			"uname -a"},
	},
	{
		client: openssh.Client{
			Host:       "some-host",
			ConfigFile: "some/config/file",
		},
		cmd: "uname -a",
		result: []string{"ssh", "-F", "some/config/file", "some-host",
			"uname -a"},
	},
	// Check for errors
	{
		client:    openssh.Client{},
		cmd:       "uname -a",
		result:    nil,
		errString: "no host given",
	},
	{
		client:    openssh.Client{Host: "some-host"},
		cmd:       "",
		result:    nil,
		errString: "no command given",
	},
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		result, err := test.client.Command(test.cmd)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.result, result)
	}
}
