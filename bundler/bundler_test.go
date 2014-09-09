package bundler_test

import (
	"os"
	"testing"

	. "github.com/mlafeldt/chef-runner/bundler"
	"github.com/stretchr/testify/assert"
)

func withPath(path string, f func()) {
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", path)
	defer os.Setenv("PATH", oldPath)
	f()
}

func TestCommand_WithoutBundlerAndGemfile(t *testing.T) {
	withPath("", func() {
		cmd := Command([]string{"rake", "test"})
		assert.Equal(t, []string{"rake", "test"}, cmd)
	})
}

func TestCommand_WithBundler(t *testing.T) {
	withPath("../testdata/bin", func() {
		cmd := Command([]string{"rake", "test"})
		assert.Equal(t, []string{"rake", "test"}, cmd)
	})
}

func TestCommand_WithGemfile(t *testing.T) {
	withPath("", func() {
		f, _ := os.Create("Gemfile")
		f.Close()
		defer os.Remove("Gemfile")

		cmd := Command([]string{"rake", "test"})
		assert.Equal(t, []string{"rake", "test"}, cmd)
	})
}

func TestCommand_WithBundlerAndGemfile(t *testing.T) {
	withPath("../testdata/bin", func() {
		f, _ := os.Create("Gemfile")
		f.Close()
		defer os.Remove("Gemfile")

		cmd := Command([]string{"rake", "test"})
		assert.Equal(t, []string{"bundle", "exec", "rake", "test"}, cmd)
	})
}
