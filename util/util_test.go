package util_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func TestFileExist(t *testing.T) {
	filename := "some-file"
	assert.False(t, util.FileExist(filename))

	f, _ := os.Create(filename)
	f.Close()
	defer os.Remove(filename)
	assert.True(t, util.FileExist(filename))
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
		assert.Equal(t, test.out, util.BaseName(test.in, test.suffix))
	}
}
