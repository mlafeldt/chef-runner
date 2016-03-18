package resolver_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner/exec"
	. "github.com/mlafeldt/chef-runner/resolver"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

const CookbookPath = "test-cookbooks"

var lastCmd []string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = args
		return nil
	})
}

func TestAutoResolve_Berkshelf(t *testing.T) {
	lastCmd = []string{}

	defer util.TestTempDir(t)()
	ioutil.WriteFile("Berksfile", []byte{}, 0644)
	os.MkdirAll(CookbookPath, 0755)
	AutoResolve(CookbookPath)

	assert.Equal(t, []string{"ruby", "-e"}, lastCmd[:2])
	assert.True(t, strings.Contains(lastCmd[2], `require "berkshelf"`))
	assert.True(t, strings.Contains(lastCmd[2], fmt.Sprintf(`.vendor("%s")`, CookbookPath)))
}

func TestAutoResolve_Librarian(t *testing.T) {
	lastCmd = []string{}

	defer util.TestTempDir(t)()
	ioutil.WriteFile("Cheffile", []byte{}, 0644)
	os.MkdirAll(CookbookPath, 0755)

	assert.NoError(t, AutoResolve(CookbookPath))
	assert.Equal(t, []string{"librarian-chef", "install", "--path", CookbookPath}, lastCmd)
}

func TestAutoResolve_Dir(t *testing.T) {
	lastCmd = []string{}

	defer util.TestTempDir(t)()
	ioutil.WriteFile("metadata.rb", []byte(`name "cats"`), 0644)

	assert.NoError(t, AutoResolve(CookbookPath))
	assert.Equal(t, []string{"rsync", "--archive", "--delete", "--compress", "metadata.rb", CookbookPath + "/cats"}, lastCmd)
}

func TestAutoResolve_DirUpdate(t *testing.T) {
	lastCmd = []string{}

	defer util.TestTempDir(t)()
	ioutil.WriteFile("metadata.rb", []byte(`name "cats"`), 0644)
	ioutil.WriteFile("Berksfile", []byte{}, 0644)
	os.MkdirAll(CookbookPath, 0755)

	assert.NoError(t, AutoResolve(CookbookPath))
	assert.Equal(t, []string{"rsync", "--archive", "--delete", "--compress", "metadata.rb", CookbookPath + "/cats"}, lastCmd)
}

func TestAutoResolve_NoCookbooks(t *testing.T) {
	lastCmd = []string{}

	defer util.TestTempDir(t)()
	err := AutoResolve(CookbookPath)

	assert.EqualError(t, err, "cookbooks could not be found")
	assert.Equal(t, []string{}, lastCmd)
}

func TestResolve_Librarian(t *testing.T) {
	lastCmd = []string{}

	defer util.TestTempDir(t)()
	os.MkdirAll(CookbookPath, 0755)

	assert.NoError(t, Resolve("librarian", CookbookPath))
	assert.Equal(t, []string{"librarian-chef", "install", "--path", CookbookPath}, lastCmd)
}
