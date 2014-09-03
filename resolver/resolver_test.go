package resolver_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/resolver"
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

	// Be quiet during testing
	log.SetLevel(log.LevelWarn)
}

func inTestDir(f func()) {
	testDir, err := util.TempDir()
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(testDir)

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := os.Chdir(testDir); err != nil {
		panic(err)
	}

	f()

	if err := os.Chdir(pwd); err != nil {
		panic(err)
	}
}

func TestAutoResolve_Berkshelf(t *testing.T) {
	lastCmd = []string{}

	inTestDir(func() {
		ioutil.WriteFile("Berksfile", []byte{}, 0644)
		os.MkdirAll(CookbookPath, 0755)

		assert.NoError(t, resolver.AutoResolve(CookbookPath))
	})

	assert.Equal(t, []string{"ruby", "-e"}, lastCmd[:2])
	assert.True(t, strings.Contains(lastCmd[2], `require "berkshelf"`))
	assert.True(t, strings.Contains(lastCmd[2],
		fmt.Sprintf(`.vendor("%s")`, CookbookPath)))
}

func TestAutoResolve_Librarian(t *testing.T) {
	lastCmd = []string{}

	inTestDir(func() {
		ioutil.WriteFile("Cheffile", []byte{}, 0644)
		os.MkdirAll(CookbookPath, 0755)

		assert.NoError(t, resolver.AutoResolve(CookbookPath))
	})

	assert.Equal(t, []string{"librarian-chef", "install", "--path", CookbookPath}, lastCmd)
}

func TestAutoResolve_Dir(t *testing.T) {
	lastCmd = []string{}

	inTestDir(func() {
		ioutil.WriteFile("metadata.rb", []byte(`name "cats"`), 0644)

		assert.NoError(t, resolver.AutoResolve(CookbookPath))
	})

	assert.Equal(t, []string{"rsync", "--archive", "--delete", "--compress",
		"--verbose", "metadata.rb", CookbookPath + "/cats"}, lastCmd)
}

func TestAutoResolve_DirUpdate(t *testing.T) {
	lastCmd = []string{}

	inTestDir(func() {
		ioutil.WriteFile("metadata.rb", []byte(`name "cats"`), 0644)
		ioutil.WriteFile("Berksfile", []byte{}, 0644)
		os.MkdirAll(CookbookPath, 0755)

		assert.NoError(t, resolver.AutoResolve(CookbookPath))
	})

	assert.Equal(t, []string{"rsync", "--archive", "--delete", "--compress",
		"--verbose", "metadata.rb", CookbookPath + "/cats"}, lastCmd)
}

func TestAutoResolve_NoCookbooks(t *testing.T) {
	lastCmd = []string{}

	inTestDir(func() {
		assert.EqualError(t, resolver.AutoResolve(CookbookPath),
			"cookbooks could not be found")
	})

	assert.Equal(t, []string{}, lastCmd)
}
