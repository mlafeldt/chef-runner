package berkshelf_test

import (
	"os"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner/resolver/berkshelf"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd := berkshelf.Command("a/b/c")
	assert.Equal(t, []string{"ruby", "-e"}, cmd[:2])
	assert.True(t, strings.Contains(cmd[2], `require "berkshelf"`))
	assert.True(t, strings.Contains(cmd[2], `.vendor("a/b/c")`))
}

func TestCommand_Bundler(t *testing.T) {
	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	cmd := berkshelf.Command("a/b/c")
	assert.Equal(t, []string{"bundle", "exec", "ruby", "-e"}, cmd[:4])
	assert.True(t, strings.Contains(cmd[4], `require "berkshelf"`))
	assert.True(t, strings.Contains(cmd[4], `.vendor("a/b/c")`))
}
