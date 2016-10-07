package util_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	. "github.com/mlafeldt/chef-runner/util"
)

func TestFileExist(t *testing.T) {
	assert.False(t, FileExist("some-non-existing-file"))
	assert.True(t, FileExist("util_test.go"))
}

func TestBaseName(t *testing.T) {
	tests := []struct {
		in     string
		suffix string
		out    string
	}{
		{"", "", "."},
		{"a", "", "a"},
		{"a/b", "", "b"},
		{"/a/b/c", "", "c"},
		{"a.x", ".x", "a"},
		{"a/b.x", ".x", "b"},
		{"a/b.x", ".y", "b.x"},
	}
	for _, test := range tests {
		assert.Equal(t, test.out, BaseName(test.in, test.suffix))
	}
}
