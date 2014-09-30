package chefsolo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/provisioner"
	. "github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Don't run any commands
	exec.SetRunnerFunc(func(args []string) error { return nil })
}

func TestProvisionerInterface(t *testing.T) {
	assert.Implements(t, (*provisioner.Provisioner)(nil), new(Provisioner))
}

// Note: Setup of cookbook dependencies is tested in the resolver package.
func TestCreateSandbox(t *testing.T) {
	util.InTestDir(func() {
		ioutil.WriteFile("Berksfile", []byte{}, 0644)
		os.MkdirAll(".chef-runner/sandbox/cookbooks", 0755)

		assert.NoError(t, Provisioner{}.CreateSandbox())

		json, err := ioutil.ReadFile(".chef-runner/sandbox/dna.json")
		assert.NoError(t, err)
		assert.Equal(t, "{}\n", string(json))

		expect := "cookbook_path \"/tmp/chef-runner/cookbooks\"\n"
		expect += "ssl_verify_mode :verify_peer\n"

		config, err := ioutil.ReadFile(".chef-runner/sandbox/solo.rb")
		assert.NoError(t, err)
		assert.Equal(t, expect, string(config))
	})
}

func TestCreateSandbox_CustomJSON(t *testing.T) {
	util.InTestDir(func() {
		ioutil.WriteFile("Berksfile", []byte{}, 0644)
		os.MkdirAll(".chef-runner/sandbox/cookbooks", 0755)

		p := Provisioner{Attributes: `{"foo": "bar"}`}
		assert.NoError(t, p.CreateSandbox())

		json, err := ioutil.ReadFile(".chef-runner/sandbox/dna.json")
		assert.NoError(t, err)
		assert.Equal(t, `{"foo": "bar"}`, string(json))

		expect := "cookbook_path \"/tmp/chef-runner/cookbooks\"\n"
		expect += "ssl_verify_mode :verify_peer\n"

		config, err := ioutil.ReadFile(".chef-runner/sandbox/solo.rb")
		assert.NoError(t, err)
		assert.Equal(t, expect, string(config))
	})
}

var provisionCommandTests = []struct {
	provisioner Provisioner
	cmd         []string
}{
	{
		Provisioner{},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "info",
		},
	},
	{
		Provisioner{
			RunList: []string{"cats::foo"},
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "info",
			"--override-runlist", "cats::foo",
		},
	},
	{
		Provisioner{
			RunList: []string{"cats::foo"},
			Format:  "null",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "null", "--log_level", "info",
			"--override-runlist", "cats::foo",
		},
	},
	{
		Provisioner{
			RunList:  []string{"cats::foo"},
			LogLevel: "error",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "error",
			"--override-runlist", "cats::foo",
		},
	},
	{
		Provisioner{
			RunList:  []string{"cats::foo", "dogs::bar"},
			Format:   "min",
			LogLevel: "warn",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "min", "--log_level", "warn",
			"--override-runlist", "cats::foo,dogs::bar",
		},
	},
	{
		Provisioner{
			RunList: []string{"cats::foo"},
			UseSudo: true,
		},
		[]string{
			"sudo", "chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "info",
			"--override-runlist", "cats::foo",
		},
	},
}

func TestProvisionCommand(t *testing.T) {
	for _, test := range provisionCommandTests {
		assert.Equal(t, test.cmd, test.provisioner.ProvisionCommand())
	}
}
