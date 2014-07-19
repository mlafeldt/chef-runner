package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRunList(t *testing.T) {
	tests := []struct {
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
	for _, test := range tests {
		runList := buildRunList(test.cookbookName, test.recipes)
		assert.Equal(t, test.runList, runList)
	}
}
