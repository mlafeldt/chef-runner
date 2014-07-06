package main

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/stretchr/testify/assert"
)

var buildRunListTests = []struct {
	cookbookName string
	recipes      []string
	runlist      string
}{
	{"cats", []string{}, "cats::default"},
	{"cats", []string{"recipes/foo.rb"}, "cats::foo"},
	{"cats", []string{"./recipes//foo.rb"}, "cats::foo"},
	{"cats", []string{"foo"}, "cats::foo"},
	{"cats", []string{"dogs::bar"}, "dogs::bar"},
	{"cats", []string{"recipes/foo.rb", "bar", "dogs::baz"}, "cats::foo,cats::bar,dogs::baz"},
}

func TestBuildRunList(t *testing.T) {
	for _, test := range buildRunListTests {
		expected := test.runlist
		actual := buildRunList(test.cookbookName, test.recipes)
		assert.Equal(t, expected, actual)
	}
}

var lastCmd string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = strings.Join(args, " ")
		return nil
	})
}

func TestOpenSSH(t *testing.T) {
	err := openSSH("somehost.local", "uname -a")
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh somehost.local -c uname -a", lastCmd)
	}
}
