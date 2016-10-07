package omnibus_test

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	. "github.com/mlafeldt/chef-runner/chef/omnibus"
	"github.com/mlafeldt/chef-runner/util"
)

func TestPrepareFiles(t *testing.T) {
	defer util.TestTempDir(t)()
	wd, _ := os.Getwd()
	i := Installer{ChefVersion: "1.2.3", SandboxPath: wd}
	assert.NoError(t, i.PrepareFiles())

	assert.True(t, util.FileExist("install.sh"))
	assert.True(t, util.FileExist("install-wrapper.sh"))
}

func TestCommand(t *testing.T) {
	tests := map[string][]string{
		"":       []string{},
		"false":  []string{},
		"latest": []string{"sudo", "sh", "/some/path/install-wrapper.sh", "/some/path/install.sh", "latest"},
		"true":   []string{"sudo", "sh", "/some/path/install-wrapper.sh", "/some/path/install.sh", "true"},
		"1.2.3":  []string{"sudo", "sh", "/some/path/install-wrapper.sh", "/some/path/install.sh", "1.2.3"},
	}
	for version, cmd := range tests {
		i := Installer{
			ChefVersion: version,
			RootPath:    "/some/path",
			Sudo:        true,
		}
		assert.Equal(t, cmd, i.Command())
	}
}
