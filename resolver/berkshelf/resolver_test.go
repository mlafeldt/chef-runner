package berkshelf_test

import (
	"strings"
	"testing"

	. "github.com/mlafeldt/chef-runner/resolver/berkshelf"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd := Command("a/b/c")
	assert.Equal(t, []string{"ruby", "-e"}, cmd[:2])
	assert.True(t, strings.Contains(cmd[2], `require "berkshelf"`))
	assert.True(t, strings.Contains(cmd[2], `.vendor("a/b/c")`))
}
