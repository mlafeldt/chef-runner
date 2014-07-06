package cookbook_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner.go/cookbook"
	"github.com/stretchr/testify/assert"
)

func TestFiles(t *testing.T) {
	expect := []string{
		"testdata/README.md",
		"testdata/metadata.rb",
		"testdata/attributes",
		"testdata/recipes",
	}
	actual, err := cookbook.Files("testdata")
	if assert.NoError(t, err) {
		assert.Equal(t, expect, actual)
	}
}

var nameFromPathTests = []struct {
	in, out string
}{
	{"/path/to/chef-cats", "cats"},
	{"/path/to/dogs-cookbook", "dogs"},
	{"some-other-name", "some-other-name"},
}

func TestNameFromPath(t *testing.T) {
	for _, test := range nameFromPathTests {
		assert.Equal(t, test.out, cookbook.NameFromPath(test.in))
	}
}
