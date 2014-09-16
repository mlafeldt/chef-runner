package omnibus_test

import (
	"testing"

	. "github.com/mlafeldt/chef-runner/omnibus"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	tests := map[string][]string{
		"":       []string{},
		"false":  []string{},
		"latest": []string{"sudo", "sh", "/some/path/install-wrapper.sh", "/some/path/install.sh", "latest"},
		"true":   []string{"sudo", "sh", "/some/path/install-wrapper.sh", "/some/path/install.sh", "true"},
		"1.2.3":  []string{"sudo", "sh", "/some/path/install-wrapper.sh", "/some/path/install.sh", "1.2.3"},
	}
	for version, cmd := range tests {
		i := Installer{ChefVersion: version, ScriptPath: "/some/path"}
		assert.Equal(t, cmd, i.Command())
	}
}
