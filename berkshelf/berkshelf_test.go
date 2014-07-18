package berkshelf_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/berkshelf"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommand(t *testing.T) {
	expect := []string{"berks", "install", "--path", "a/b/c"}
	actual := berkshelf.InstallCommand("a/b/c")
	assert.Equal(t, expect, actual)
}

func TestInstallCommand_Bundler(t *testing.T) {
	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	expect := []string{"bundle", "exec", "berks", "install", "--path", "a/b/c"}
	actual := berkshelf.InstallCommand("a/b/c")
	assert.Equal(t, expect, actual)
}
