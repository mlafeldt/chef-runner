package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
