package berkshelf_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/berkshelf"
	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/stretchr/testify/assert"
)

var lastCmd string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = strings.Join(args, " ")
		return nil
	})
}

func TestInstall(t *testing.T) {
	err := berkshelf.Install("a/b/c")
	if assert.NoError(t, err) {
		assert.Equal(t, "berks install --path a/b/c", lastCmd)
	}
}

func TestInstall_Bundler(t *testing.T) {
	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	err := berkshelf.Install("a/b/c")
	if assert.NoError(t, err) {
		assert.Equal(t, "bundle exec berks install --path a/b/c", lastCmd)
	}
}