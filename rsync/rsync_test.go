package rsync_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/rsync"
	"github.com/stretchr/testify/assert"
)

var lastCmd string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = strings.Join(args, " ")
		return nil
	})
}

var copyTests = []struct {
	client rsync.Client
	src    []string
	dst    string
	cmd    string
}{
	{
		rsync.Client{},
		[]string{"a"}, "b",
		"rsync a b",
	},
	{
		rsync.Client{},
		[]string{"a", "b"}, "c",
		"rsync a b c",
	},
	{
		rsync.Client{Archive: true},
		[]string{"a"}, "b",
		"rsync --archive a b",
	},
	{
		rsync.Client{Delete: true},
		[]string{"a"}, "b",
		"rsync --delete a b",
	},
	{
		rsync.Client{Verbose: true},
		[]string{"a"}, "b",
		"rsync --verbose a b",
	},
	{
		rsync.Client{Exclude: []string{"x", "y"}},
		[]string{"a"}, "b",
		"rsync --exclude x --exclude y a b",
	},
	{
		rsync.Client{Archive: true, Delete: true, Exclude: []string{"x"}},
		[]string{"a"}, "b",
		"rsync --archive --delete --exclude x a b",
	},
}

func TestCopy(t *testing.T) {
	for _, test := range copyTests {
		err := test.client.Copy(test.src, test.dst)
		if assert.NoError(t, err) {
			assert.Equal(t, test.cmd, lastCmd)
		}
	}
}

func TestCopy_MissingSource(t *testing.T) {
	err := rsync.DefaultClient.Copy([]string{}, "a/b")
	assert.EqualError(t, err, "No source given")
}

func TestCopy_MissingDestination(t *testing.T) {
	err := rsync.DefaultClient.Copy([]string{"a"}, "")
	assert.EqualError(t, err, "No destination given")
}
