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
	c, err := openssh.NewClient("some-host")
	assert.NoError(t, err)
	assert.Equal(t, expect, c.String())
}

var commandTests = []struct {
	client openssh.Client
	cmd    []string
}{
	{
		openssh.Client{
			Host: "some-host",
		},
		[]string{"ssh", "some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host: "some-host",
			User: "some-user",
		},
		[]string{"ssh", "-l", "some-user", "some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host: "some-host",
			Port: 1234,
		},
		[]string{"ssh", "-p", "1234", "some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host:        "some-host",
			PrivateKeys: []string{"some-key", "another-key"},
		},
		[]string{"ssh", "-i", "some-key", "-i", "another-key",
			"some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host: "some-host",
			Options: map[string]string{
				"SomeOption":    "yes",
				"AnotherOption": "no",
			},
		},
		[]string{"ssh", "-o", "AnotherOption=no", "-o", "SomeOption=yes",
			"some-host", "uname -a"},
	},
	{
		openssh.Client{
			Host:        "some-host",
			User:        "some-user",
			Port:        1234,
			PrivateKeys: []string{"some-key"},
			Options:     map[string]string{"SomeOption": "yes"},
		},
		[]string{"ssh", "-l", "some-user", "-p", "1234", "-i", "some-key",
			"-o", "SomeOption=yes", "some-host", "uname -a"},
	},
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		cmd, err := test.client.Command("uname -a")
		if assert.NoError(t, err) {
			assert.Equal(t, test.cmd, cmd)
		}
	}
}

func TestRunCommand_MissingCommand(t *testing.T) {
	err := openssh.Client{Host: "some-host"}.RunCommand("")
	assert.EqualError(t, err, "no command given")
}

func TestRunCommand_MissingHost(t *testing.T) {
	err := openssh.Client{}.RunCommand("uname -a")
	assert.EqualError(t, err, "no host given")
}
