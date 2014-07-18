package provisioner_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/provisioner"
	"github.com/stretchr/testify/assert"
)

var sandboxPathToTests = []struct {
	in, out string
}{
	{"", ".chef-runner"},
	{"a", ".chef-runner/a"},
	{"/a/b/c", ".chef-runner/a/b/c"},
}

func TestSandboxPathTo(t *testing.T) {
	for _, test := range sandboxPathToTests {
		assert.Equal(t, test.out, provisioner.SandboxPathTo(test.in))
	}
}

var rootPathToTests = []struct {
	in, out string
}{
	{"", "/vagrant/.chef-runner"},
	{"a", "/vagrant/.chef-runner/a"},
	{"/a/b/c", "/vagrant/.chef-runner/a/b/c"},
}

func TestRootPathTo(t *testing.T) {
	for _, test := range rootPathToTests {
		assert.Equal(t, test.out, provisioner.RootPathTo(test.in))
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
