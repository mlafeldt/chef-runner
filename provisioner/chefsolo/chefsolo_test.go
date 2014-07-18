package chefsolo_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/stretchr/testify/assert"
)

var cmdPrefix = []string{"sudo", "chef-solo",
	"--config", "/vagrant/.chef-runner/solo.rb",
	"--json-attributes", "/vagrant/.chef-runner/dna.json",
}

var commandTests = []struct {
	provisioner chefsolo.Provisoner
	cmd         []string
}{
	{
		chefsolo.Provisoner{
			RunList: []string{"cats::foo"},
		},
		append(cmdPrefix, "--override-runlist", "cats::foo",
			"--format", "null", "--log_level", "info"),
	},
	{
		chefsolo.Provisoner{
			RunList: []string{"cats::foo"},
			Format:  "doc",
		},
		append(cmdPrefix, "--override-runlist", "cats::foo",
			"--format", "doc", "--log_level", "info"),
	},
	{
		chefsolo.Provisoner{
			RunList:  []string{"cats::foo"},
			LogLevel: "error",
		},
		append(cmdPrefix, "--override-runlist", "cats::foo",
			"--format", "null", "--log_level", "error"),
	},
	{
		chefsolo.Provisoner{
			RunList:  []string{"cats::foo", "dogs::bar"},
			Format:   "min",
			LogLevel: "warn",
		},
		append(cmdPrefix, "--override-runlist", "cats::foo,dogs::bar",
			"--format", "min", "--log_level", "warn"),
	},
}

func TestCommand(t *testing.T) {
	for _, test := range commandTests {
		assert.Equal(t, test.cmd, test.provisioner.Command())
	}
}
