package cookbook_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner.go/cookbook"
	"github.com/stretchr/testify/assert"
)

func TestNewCookbook(t *testing.T) {
	cb, err := cookbook.NewCookbook("testdata")
	assert.NoError(t, err)
	if assert.NotNil(t, cb) {
		assert.Equal(t, "testdata", cb.Path)
		assert.Equal(t, "practicingruby", cb.Name)
		assert.Equal(t, "1.3.1", cb.Version)
	}
}

func TestFiles(t *testing.T) {
	cb, _ := cookbook.NewCookbook("testdata")
	expect := []string{
		"testdata/README.md",
		"testdata/metadata.rb",
		"testdata/attributes",
		"testdata/recipes",
	}
	actual, err := cb.Files()
	if assert.NoError(t, err) {
		assert.Equal(t, expect, actual)
	}
}
