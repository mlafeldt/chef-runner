package bundler_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/bundler"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	cmd := bundler.Command([]string{"rake", "test"})
	assert.Equal(t, []string{"rake", "test"}, cmd)
}

func TestCommand_WithGemfile(t *testing.T) {
	f, _ := os.Create("Gemfile")
	f.Close()
	defer os.Remove("Gemfile")

	cmd := bundler.Command([]string{"rake", "test"})
	assert.Equal(t, []string{"bundle", "exec", "rake", "test"}, cmd)
}
