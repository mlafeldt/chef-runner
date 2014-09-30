package provisioner_test

import (
	"os"
	"testing"

	. "github.com/mlafeldt/chef-runner/provisioner"
	"github.com/stretchr/testify/assert"
)

func TestSandboxPathTo(t *testing.T) {
	tests := map[string]string{
		"":       ".chef-runner/sandbox",
		"a":      ".chef-runner/sandbox/a",
		"/a/b/c": ".chef-runner/sandbox/a/b/c",
	}
	for in, out := range tests {
		assert.Equal(t, out, SandboxPathTo(in))
	}
	assert.Equal(t, ".chef-runner/sandbox/a/b/c",
		SandboxPathTo("a", "b", "c"))
}

func TestRootPathTo(t *testing.T) {
	tests := map[string]string{
		"":       "/tmp/chef-runner",
		"a":      "/tmp/chef-runner/a",
		"/a/b/c": "/tmp/chef-runner/a/b/c",
	}
	for in, out := range tests {
		assert.Equal(t, out, RootPathTo(in))
	}
	assert.Equal(t, "/tmp/chef-runner/a/b/c",
		RootPathTo("a", "b", "c"))
}

func TestCreateAndCleanupSandbox(t *testing.T) {
	err := CreateSandbox()
	if assert.NoError(t, err) {
		fi, err := os.Stat(SandboxPath)
		if assert.NoError(t, err) {
			assert.True(t, fi.IsDir())
		}
	}

	err = CleanupSandbox()
	if assert.NoError(t, err) {
		_, err := os.Stat(SandboxPath)
		assert.Error(t, err)
	}
}
