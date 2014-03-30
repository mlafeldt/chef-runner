package main

import (
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/stretchr/testify/assert"
)

var cookbookNameTests = []struct {
	cookbookPath, cookbookName string
}{
	{"/path/to/chef-cats", "cats"},
	{"/path/to/dogs-cookbook", "dogs"},
	{"some-other-name", "some-other-name"},
}

func TestCookbookName(t *testing.T) {
	for _, test := range cookbookNameTests {
		expected := test.cookbookName
		actual := cookbookName(test.cookbookPath)
		assert.Equal(t, expected, actual)
	}
}

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

var history []string

func clearHistory() { history = []string{} }

func init() {
	f := func(args []string) error {
		history = append(history, strings.Join(args, " "))
		return nil
	}
	exec.SetRunnerFunc(f)
}

func TestOpenSSH(t *testing.T) {
	defer clearHistory()

	openSSH("somehost.local", "uname -a")

	if assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "ssh somehost.local -c uname -a", history[0])
	}
}
