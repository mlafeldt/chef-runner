package openssh_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/stretchr/testify/assert"
)

var newClientTests = map[string]*openssh.Client{
	"some-host":                &openssh.Client{Host: "some-host"},
	"some-user@some-host":      &openssh.Client{Host: "some-host", User: "some-user"},
	"some-host:1234":           &openssh.Client{Host: "some-host", Port: 1234},
	"some-user@some-host:1234": &openssh.Client{Host: "some-host", User: "some-user", Port: 1234},
}

func TestNewClient(t *testing.T) {
	for host, client := range newClientTests {
		result, err := openssh.NewClient(host)
		assert.NoError(t, err)
		assert.Equal(t, client, result)
	}
}

func TestNewClient_InvalidPort(t *testing.T) {
	c, err := openssh.NewClient("some-host:abc")
	assert.EqualError(t, err, "invalid SSH port")
	assert.Nil(t, c)
}

var commandTests = []struct {
	client openssh.Client
	cmd    string
	result []string
}{
	{
		client: openssh.Client{},
		cmd:    "",
		result: []string{"ssh"},
	},
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
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		result := test.client.Command(test.cmd)
		assert.Equal(t, test.result, result)
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

func TestShell(t *testing.T) {
	c, _ := openssh.NewClient("some-user@some-host:1234")
	assert.Equal(t, []string{"ssh", "-l", "some-user", "-p", "1234"}, c.Shell())
}
