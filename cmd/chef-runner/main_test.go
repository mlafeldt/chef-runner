package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var buildRunListTests = []struct {
	cookbookName string
	recipes      []string
	runList      []string
}{
	{"cats", []string{}, []string{"cats::default"}},
	{"cats", []string{"recipes/foo.rb"}, []string{"cats::foo"}},
	{"cats", []string{"./recipes//foo.rb"}, []string{"cats::foo"}},
	{"cats", []string{"foo"}, []string{"cats::foo"}},
	{"cats", []string{"dogs::bar"}, []string{"dogs::bar"}},
	{"cats", []string{"recipes/foo.rb", "bar", "dogs::baz"},
		[]string{"cats::foo", "cats::bar", "dogs::baz"}},
}

func TestBuildRunList(t *testing.T) {
	for _, test := range buildRunListTests {
		runList := buildRunList(test.cookbookName, test.recipes)
		assert.Equal(t, test.runList, runList)
	}
}
