package librarian_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/librarian"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommand(t *testing.T) {
	expect := []string{"librarian-chef", "install", "--path", "a/b/c"}
	actual := librarian.InstallCommand("a/b/c")
	assert.Equal(t, expect, actual)
}

func TestInstallCommand_Bundler(t *testing.T) {
	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	expect := []string{"bundle", "exec", "librarian-chef", "install", "--path", "a/b/c"}
	actual := librarian.InstallCommand("a/b/c")
	assert.Equal(t, expect, actual)
}
