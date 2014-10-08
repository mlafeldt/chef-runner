package openssh_test

import (
	"testing"

	. "github.com/mlafeldt/chef-runner/openssh"
	"github.com/stretchr/testify/assert"
)

var newClientTests = map[string]*Client{
	"some-host":                &Client{Host: "some-host"},
	"some-user@some-host":      &Client{Host: "some-host", User: "some-user"},
	"some-host:1234":           &Client{Host: "some-host", Port: 1234},
	"some-user@some-host:1234": &Client{Host: "some-host", User: "some-user", Port: 1234},
}

func TestNewClient(t *testing.T) {
	for host, client := range newClientTests {
		result, err := NewClient(host)
		assert.NoError(t, err)
		assert.Equal(t, client, result)
	}
}

func TestNewClient_InvalidPort(t *testing.T) {
	c, err := NewClient("some-host:abc")
	assert.EqualError(t, err, "invalid SSH port")
	assert.Nil(t, c)
}

var commandTests = []struct {
	client Client
	args   []string
	result []string
}{
	{
		client: Client{},
		args:   []string{},
		result: []string{"ssh"},
	},
	{
		client: Client{
			Host: "some-host",
		},
		args:   []string{"uname", "-a"},
		result: []string{"ssh", "some-host", "uname", "-a"},
	},
	{
		client: Client{
			Host: "some-host",
			User: "some-user",
		},
		args:   []string{"uname", "-a"},
		result: []string{"ssh", "-l", "some-user", "some-host", "uname", "-a"},
	},
	{
		client: Client{
			Host: "some-host",
			Port: 1234,
		},
		args:   []string{"uname", "-a"},
		result: []string{"ssh", "-p", "1234", "some-host", "uname", "-a"},
	},
	{
		client: Client{
			Host:        "some-host",
			PrivateKeys: []string{"some-key", "another-key"},
		},
		args: []string{"uname", "-a"},
		result: []string{"ssh", "-i", "some-key", "-i", "another-key",
			"some-host", "uname", "-a"},
	},
	{
		client: Client{
			Host: "some-host",
			Options: []string{
				"SomeOption=yes",
				"AnotherOption 1 2 3",
			},
		},
		args: []string{"uname", "-a"},
		result: []string{"ssh", "-o", "SomeOption=yes", "-o", "AnotherOption 1 2 3",
			"some-host", "uname", "-a"},
	},
	{
		client: Client{
			Host:        "some-host",
			User:        "some-user",
			Port:        1234,
			PrivateKeys: []string{"some-key"},
			Options:     []string{"SomeOption=yes"},
		},
		args: []string{"uname", "-a"},
		result: []string{"ssh", "-l", "some-user", "-p", "1234",
			"-i", "some-key", "-o", "SomeOption=yes", "some-host",
			"uname", "-a"},
	},
	{
		client: Client{
			Host:       "some-host",
			ConfigFile: "some/config/file",
		},
		args: []string{"uname", "-a"},
		result: []string{"ssh", "-F", "some/config/file", "some-host",
			"uname", "-a"},
	},
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		result := test.client.Command(test.args)
		assert.Equal(t, test.result, result)
	}
}

func TestRunCommand_MissingCommand(t *testing.T) {
	err := Client{Host: "some-host"}.RunCommand([]string{})
	assert.EqualError(t, err, "no command given")
}

func TestRunCommand_MissingHost(t *testing.T) {
	err := Client{}.RunCommand([]string{"uname", "-a"})
	assert.EqualError(t, err, "no host given")
}

var shellTests = []struct {
	client Client
	shell  string
}{
	{
		Client{Host: "some-host", User: "some-user", Port: 1234},
		`"ssh" "-l" "some-user" "-p" "1234"`,
	},
	{
		Client{Host: "some-host", Options: []string{"x=1"}},
		`"ssh" "-o" "x=1"`,
	},
	{
		Client{Host: "some-host", Options: []string{"y 2 3"}},
		`"ssh" "-o" "y 2 3"`,
	},
}

func TestShell(t *testing.T) {
	for _, test := range shellTests {
		assert.Equal(t, test.shell, test.client.Shell())
	}
}
