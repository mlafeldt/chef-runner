package provisioner_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/provisioner"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Be quiet during testing
	log.SetLevel(log.LevelWarn)
}

func TestSandboxPathTo(t *testing.T) {
	tests := map[string]string{
		"":       ".chef-runner",
		"a":      ".chef-runner/a",
		"/a/b/c": ".chef-runner/a/b/c",
	}
	for in, out := range tests {
		assert.Equal(t, out, provisioner.SandboxPathTo(in))
	}
}

func TestRootPathTo(t *testing.T) {
	tests := map[string]string{
		"":       "/vagrant/.chef-runner",
		"a":      "/vagrant/.chef-runner/a",
		"/a/b/c": "/vagrant/.chef-runner/a/b/c",
	}
	for in, out := range tests {
		assert.Equal(t, out, provisioner.RootPathTo(in))
	}
}

func TestCreateAndCleanupSandbox(t *testing.T) {
	err := provisioner.CreateSandbox()
	if assert.NoError(t, err) {
		fi, err := os.Stat(provisioner.SandboxPath)
		if assert.NoError(t, err) {
			assert.True(t, fi.IsDir())
		}
	}

	err = provisioner.CleanupSandbox()
	if assert.NoError(t, err) {
		_, err := os.Stat(provisioner.SandboxPath)
		assert.Error(t, err)
	}
}
