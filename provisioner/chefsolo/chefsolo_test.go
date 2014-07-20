package chefsolo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/stretchr/testify/assert"
)

var lastCmd []string

func init() {
	exec.SetRunnerFunc(func(args []string) error {
		lastCmd = args
		return nil
	})

	// Be quiet during testing
	log.SetLevel(log.LevelWarn)
}

var createSandboxTests = []struct {
	provisioner   chefsolo.Provisoner
	fakeCookbooks bool

	writeAttributes string
	writeConfig     string
	runCmd          []string
}{
	{
		chefsolo.Provisoner{},
		false,
		"{}\n",
		"cookbook_path \"/vagrant/.chef-runner/cookbooks\"\n",
		[]string{"bundle", "exec", "berks", "install",
			"--path", ".chef-runner/cookbooks"},
	},
	{
		chefsolo.Provisoner{},
		true,
		"{}\n",
		"cookbook_path \"/vagrant/.chef-runner/cookbooks\"\n",
		[]string{"rsync", "--archive", "--delete", "--verbose",
			"README.md", "metadata.rb", "attributes", "recipes",
			".chef-runner/cookbooks/practicingruby"},
	},
	{
		chefsolo.Provisoner{Attributes: `{"foo": "bar"}`},
		false,
		`{"foo": "bar"}`,
		"cookbook_path \"/vagrant/.chef-runner/cookbooks\"\n",
		[]string{"bundle", "exec", "berks", "install",
			"--path", ".chef-runner/cookbooks"},
	},
}

func TestCreateSandbox(t *testing.T) {
	if err := os.Chdir("../../testdata"); err != nil {
		panic(err)
	}

	defer os.RemoveAll(".chef-runner")

	for _, test := range createSandboxTests {
		if test.fakeCookbooks {
			os.MkdirAll(".chef-runner/cookbooks", 0755)
		}

		assert.NoError(t, test.provisioner.CreateSandbox())

		attributes, err := ioutil.ReadFile(".chef-runner/dna.json")
		assert.NoError(t, err)
		assert.Equal(t, test.writeAttributes, string(attributes))

		config, err := ioutil.ReadFile(".chef-runner/solo.rb")
		assert.NoError(t, err)
		assert.Equal(t, test.writeConfig, string(config))

		assert.Equal(t, test.runCmd, lastCmd)

		assert.NoError(t, test.provisioner.CleanupSandbox())
	}
}

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
