package rsync_test

import (
	"testing"

	. "github.com/mlafeldt/chef-runner/rsync"
	"github.com/stretchr/testify/assert"
)

var commandTests = []struct {
	client    Client
	src       []string
	dst       string
	cmd       []string
	errString string
}{
	{
		client: Client{},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "a", "b"},
	},
	{
		client: Client{},
		src:    []string{"a", "b"},
		dst:    "c",
		cmd:    []string{"rsync", "a", "b", "c"},
	},
	{
		client: Client{Archive: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--archive", "a", "b"},
	},
	{
		client: Client{Delete: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--delete", "a", "b"},
	},
	{
		client: Client{Compress: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--compress", "a", "b"},
	},
	{
		client: Client{Verbose: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--verbose", "a", "b"},
	},
	{
		client: Client{Exclude: []string{"x", "y"}},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--exclude", "x", "--exclude", "y", "a", "b"},
	},
	{
		client: Client{RemoteShell: "some-shell", RemoteHost: "some-host"},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--rsh", "some-shell", "a", "some-host:b"},
	},
	{
		client: Client{Archive: true, Compress: true, Exclude: []string{"x"}},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--archive", "--compress", "--exclude", "x", "a", "b"},
	},
	// Check for errors
	{
		client:    Client{},
		src:       []string{},
		dst:       "b",
		cmd:       nil,
		errString: "no source given",
	},
	{
		client:    Client{},
		src:       []string{"a"},
		dst:       "",
		cmd:       nil,
		errString: "no destination given",
	},
	{
		client:    Client{RemoteShell: "some-shell"},
		src:       []string{"a"},
		dst:       "b",
		cmd:       nil,
		errString: "no remote host given",
	},
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		cmd, err := test.client.Command(test.dst, test.src...)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.cmd, cmd)
	}
}
