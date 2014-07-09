package main

import (
	"testing"

	"github.com/mlafeldt/chef-runner.go/cookbook"
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
		cb := cookbook.Cookbook{Name: test.cookbookName}
		assert.Equal(t, test.runlist, buildRunList(&cb, test.recipes))
	}
}
