package cookbook_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/mlafeldt/chef-runner/chef/cookbook"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func TestNewCookbook(t *testing.T) {
	cb, err := NewCookbook("../../testdata")
	assert.NoError(t, err)
	if assert.NotNil(t, cb) {
		assert.Equal(t, "../../testdata", cb.Path)
		assert.Equal(t, "practicingruby", cb.Name)
		assert.Equal(t, "1.3.1", cb.Version)
	}
}

func TestNewCookbook_WithoutMetadata(t *testing.T) {
	cb, err := NewCookbook(".")
	assert.NoError(t, err)
	if assert.NotNil(t, cb) {
		assert.Equal(t, ".", cb.Path)
		assert.Equal(t, "", cb.Name)
		assert.Equal(t, "", cb.Version)
	}
}

func TestString(t *testing.T) {
	cb := Cookbook{Name: "cats", Version: "1.2.3"}
	assert.Equal(t, "cats 1.2.3", cb.String())
}

func TestFiles(t *testing.T) {
	cb, _ := NewCookbook("../../testdata")
	expect := []string{
		"../../testdata/README.md",
		"../../testdata/metadata.rb",
		"../../testdata/attributes",
		"../../testdata/recipes",
	}
	assert.Equal(t, expect, cb.Files())
}

func TestStrip(t *testing.T) {
	util.InTestDir(func() {
		for _, f := range []string{"CHANGELOG.md", "README.md", "metadata.rb"} {
			ioutil.WriteFile(f, []byte{}, 0644)
		}
		for _, d := range []string{"attributes", "recipes", "tmp"} {
			os.Mkdir(d, 0755)
		}

		cb, _ := NewCookbook(".")
		assert.NoError(t, cb.Strip())

		expect := []string{
			"README.md",
			"attributes",
			"metadata.rb",
			"recipes",
		}

		files, err := ioutil.ReadDir(".")
		if err != nil {
			panic(err)
		}

		var actual []string
		for _, f := range files {
			actual = append(actual, f.Name())
		}

		assert.Equal(t, expect, actual)
	})
}
