package chefsolo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Don't run any commands
	exec.SetRunnerFunc(func(args []string) error { return nil })

	// Be quiet during testing
	log.SetLevel(log.LevelWarn)
}

func inTestDir(f func()) {
	testDir, err := util.TempDir()
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(testDir)

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := os.Chdir(testDir); err != nil {
		panic(err)
	}

	f()

	if err := os.Chdir(pwd); err != nil {
		panic(err)
	}
}

// Note: Setup of cookbook dependencies is tested in the resolver package.
func TestCreateSandbox(t *testing.T) {
	inTestDir(func() {
		ioutil.WriteFile("Berksfile", []byte{}, 0644)

		assert.NoError(t, chefsolo.Provisioner{}.CreateSandbox())

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
	inTestDir(func() {
		ioutil.WriteFile("Berksfile", []byte{}, 0644)

		p := chefsolo.Provisioner{Attributes: `{"foo": "bar"}`}
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

var cmdPrefix = []string{"sudo", "chef-solo",
	"--config", "/tmp/chef-runner/solo.rb",
	"--json-attributes", "/tmp/chef-runner/dna.json",
}

var commandTests = []struct {
	provisioner chefsolo.Provisioner
	cmd         []string
}{
	{
		chefsolo.Provisioner{
			RunList: []string{"cats::foo"},
		},
		append(cmdPrefix, "--override-runlist", "cats::foo",
			"--format", "null", "--log_level", "info"),
	},
	{
		chefsolo.Provisioner{
			RunList: []string{"cats::foo"},
			Format:  "doc",
		},
		append(cmdPrefix, "--override-runlist", "cats::foo",
			"--format", "doc", "--log_level", "info"),
	},
	{
		chefsolo.Provisioner{
			RunList:  []string{"cats::foo"},
			LogLevel: "error",
		},
		append(cmdPrefix, "--override-runlist", "cats::foo",
			"--format", "null", "--log_level", "error"),
	},
	{
		chefsolo.Provisioner{
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
