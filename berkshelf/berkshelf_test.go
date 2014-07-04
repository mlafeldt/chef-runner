package berkshelf_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner.go/berkshelf"
	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/stretchr/testify/assert"
)

var history []string

func clearHistory() { history = []string{} }

func init() {
	f := func(args []string) error {
		history = append(history, strings.Join(args, " "))
		return nil
	}
	exec.SetRunnerFunc(f)
}

func TestInstall(t *testing.T) {
	defer clearHistory()

	err := berkshelf.Install("a/b/c")
	if assert.NoError(t, err) && assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "berks install --path a/b/c", history[0])
	}
}

func TestInstall_Bundler(t *testing.T) {
	defer clearHistory()

	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	err := berkshelf.Install("a/b/c")
	if assert.NoError(t, err) && assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "bundle exec berks install --path a/b/c", history[0])
	}
}
