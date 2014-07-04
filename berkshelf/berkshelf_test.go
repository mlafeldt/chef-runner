package berkshelf_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/berkshelf"
	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/stretchr/testify/assert"
)

var last_cmd string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		last_cmd = strings.Join(args, " ")
		return nil
	})
}

func TestInstall(t *testing.T) {
	err := berkshelf.Install("a/b/c")
	if assert.NoError(t, err) {
		assert.Equal(t, "berks install --path a/b/c", last_cmd)
	}
}

func TestInstall_Bundler(t *testing.T) {
	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	err := berkshelf.Install("a/b/c")
	if assert.NoError(t, err) {
		assert.Equal(t, "bundle exec berks install --path a/b/c", last_cmd)
	}
}
