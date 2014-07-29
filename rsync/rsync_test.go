package rsync_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/rsync"
	"github.com/stretchr/testify/assert"
)

var commandTests = []struct {
	client    rsync.Client
	src       []string
	dst       string
	cmd       []string
	errString string
}{
	{
		client: rsync.Client{},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "a", "b"},
	},
	{
		client: rsync.Client{},
		src:    []string{"a", "b"},
		dst:    "c",
		cmd:    []string{"rsync", "a", "b", "c"},
	},
	{
		client: rsync.Client{Archive: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--archive", "a", "b"},
	},
	{
		client: rsync.Client{Delete: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--delete", "a", "b"},
	},
	{
		client: rsync.Client{Verbose: true},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--verbose", "a", "b"},
	},
	{
		client: rsync.Client{Exclude: []string{"x", "y"}},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--exclude", "x", "--exclude", "y", "a", "b"},
	},
	{
		client: rsync.Client{Archive: true, Delete: true, Exclude: []string{"x"}},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--archive", "--delete", "--exclude", "x", "a", "b"},
	},
	{
		client: rsync.Client{RemoteShell: "some-shell", RemoteHost: "some-host"},
		src:    []string{"a"},
		dst:    "b",
		cmd:    []string{"rsync", "--rsh", "some-shell", "a", "some-host:b"},
	},
	// Check for errors
	{
		client:    rsync.Client{},
		src:       []string{},
		dst:       "b",
		cmd:       nil,
		errString: "no source given",
	},
	{
		client:    rsync.Client{},
		src:       []string{"a"},
		dst:       "",
		cmd:       nil,
		errString: "no destination given",
	},
	{
		client:    rsync.Client{RemoteShell: "some-shell"},
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
