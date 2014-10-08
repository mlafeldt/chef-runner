package cli_test

import (
	"os"
	"testing"

	. "github.com/mlafeldt/chef-runner/cli"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		args      []string
		flags     *Flags
		errString string
	}{
		{
			args:  []string{},
			flags: &Flags{},
		},
		{
			args:  []string{"-H", "some-host"},
			flags: &Flags{Host: "some-host"},
		},
		{
			args:  []string{"--host", "some-host"},
			flags: &Flags{Host: "some-host"},
		},
		{
			args:  []string{"-M", "some-machine"},
			flags: &Flags{Machine: "some-machine"},
		},
		{
			args:  []string{"--machine", "some-machine"},
			flags: &Flags{Machine: "some-machine"},
		},
		{
			args:  []string{"-K", "some-instance"},
			flags: &Flags{Kitchen: "some-instance"},
		},
		{
			args:  []string{"--kitchen", "some-instance"},
			flags: &Flags{Kitchen: "some-instance"},
		},
		{
			args:  []string{"--ssh-option", "x=1", "--ssh-option", "y 2 3"},
			flags: &Flags{SSHOptions: []string{"x=1", "y 2 3"}},
		},
		{
			args:  []string{"-i", "1.2.3"},
			flags: &Flags{ChefVersion: "1.2.3"},
		},
		{
			args:  []string{"--install-chef", "1.2.3"},
			flags: &Flags{ChefVersion: "1.2.3"},
		},
		{
			args:  []string{"-F", "some-format"},
			flags: &Flags{Format: "some-format"},
		},
		{
			args:  []string{"--format", "some-format"},
			flags: &Flags{Format: "some-format"},
		},
		{
			args:  []string{"-l", "some-level"},
			flags: &Flags{LogLevel: "some-level"},
		},
		{
			args:  []string{"--log_level", "some-level"},
			flags: &Flags{LogLevel: "some-level"},
		},
		{
			args:  []string{"-j", "some-file"},
			flags: &Flags{JSONFile: "some-file"},
		},
		{
			args:  []string{"--json-attributes", "some-file"},
			flags: &Flags{JSONFile: "some-file"},
		},
		{
			args:  []string{"--version"},
			flags: &Flags{ShowVersion: true},
		},
		{
			args:  []string{"some-recipe", "another-recipe"},
			flags: &Flags{Recipes: []string{"some-recipe", "another-recipe"}},
		},
		{
			args: []string{"--machine", "some-machine", "-l", "some-level", "-i", "true", "some-recipe"},
			flags: &Flags{
				Machine:     "some-machine",
				ChefVersion: "true",
				LogLevel:    "some-level",
				Recipes:     []string{"some-recipe"},
			},
		},
		// Check for errors
		{
			args:      []string{"-H", "some-host", "-M", "some-machine"},
			flags:     nil,
			errString: "-H, -M, and -K cannot be used together",
		},
		{
			args:      []string{"-K", "some-instance", "-M", "some-machine"},
			flags:     nil,
			errString: "-H, -M, and -K cannot be used together",
		},
	}
	for _, test := range tests {
		flags, err := ParseFlags(test.args)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.flags, flags)
	}
}

func TestLogLevel(t *testing.T) {
	tests := map[string]log.Level{
		"":      log.LevelInfo,
		"debug": log.LevelDebug,
		"info":  log.LevelInfo,
		"warn":  log.LevelWarn,
		"error": log.LevelError,
		"DEBUG": log.LevelDebug,
		"INFO":  log.LevelInfo,
		"WARN":  log.LevelWarn,
		"ERROR": log.LevelError,
		"foo":   log.LevelInfo,
	}
	defer os.Setenv("CHEF_RUNNER_LOG", "")
	for env, level := range tests {
		os.Setenv("CHEF_RUNNER_LOG", env)
		assert.Equal(t, level, LogLevel())
	}
}
