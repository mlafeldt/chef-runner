package rsync_test

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/rsync"
	"github.com/stretchr/testify/assert"
)

var history []string

func clearHistory() { history = []string{} }

func init() {
	f := func(args []string) error {
		history = append(history, strings.Join(args, " "))
		return nil
	}
	exec.SetRunnerFunc(f)
}

var copyTests = []struct {
	src  []string
	dst  string
	opts rsync.Options
	cmd  string
}{
	{
		[]string{"a"}, "b",
		rsync.Options{},
		"rsync a b",
	},
	{
		[]string{"a", "b"}, "c",
		rsync.Options{},
		"rsync a b c",
	},
	{
		[]string{"a"}, "b",
		rsync.Options{Archive: true},
		"rsync --archive a b",
	},
	{
		[]string{"a"}, "b",
		rsync.Options{Delete: true},
		"rsync --delete a b",
	},
	{
		[]string{"a"}, "b",
		rsync.Options{Verbose: true},
		"rsync --verbose a b",
	},
	{
		[]string{"a"}, "b",
		rsync.Options{Exclude: []string{"x", "y"}},
		"rsync --exclude x --exclude y a b",
	},
	{
		[]string{"a"}, "b",
		rsync.Options{Archive: true, Delete: true, Exclude: []string{"x"}},
		"rsync --archive --delete --exclude x a b",
	},
}

func TestCopy(t *testing.T) {
	for _, test := range copyTests {
		err := rsync.Copy(test.src, test.dst, test.opts)
		if assert.NoError(t, err) && assert.Equal(t, 1, len(history)) {
			assert.Equal(t, test.cmd, history[0])
		}
		clearHistory()
	}
}
