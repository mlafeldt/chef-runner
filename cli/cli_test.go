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
			flags: &Flags{Sudo: true, Color: true},
		},
		{
			args:  []string{"-H", "some-host"},
			flags: &Flags{Host: "some-host", Sudo: true, Color: true},
		},
		{
			args:  []string{"--host", "some-host"},
			flags: &Flags{Host: "some-host", Sudo: true, Color: true},
		},
		{
			args:  []string{"-M", "some-machine"},
			flags: &Flags{Machine: "some-machine", Sudo: true, Color: true},
		},
		{
			args:  []string{"--machine", "some-machine"},
			flags: &Flags{Machine: "some-machine", Sudo: true, Color: true},
		},
		{
			args:  []string{"-K", "some-instance"},
			flags: &Flags{Kitchen: "some-instance", Sudo: true, Color: true},
		},
		{
			args:  []string{"--kitchen", "some-instance"},
			flags: &Flags{Kitchen: "some-instance", Sudo: true, Color: true},
		},
		{
			args:  []string{"-L"},
			flags: &Flags{Local: true, Sudo: true, Color: true},
		},
		{
			args:  []string{"--local"},
			flags: &Flags{Local: true, Sudo: true, Color: true},
		},
		{
			args:  []string{"--ssh", "x=1", "--ssh", "y 2 3"},
			flags: &Flags{SSHOptions: []string{"x=1", "y 2 3"}, Sudo: true, Color: true},
		},
		{
			args:  []string{"--rsync", "-x", "--rsync", "--y"},
			flags: &Flags{RsyncOptions: []string{"-x", "--y"}, Sudo: true, Color: true},
		},
		{
			args:  []string{"--resolver", "some-resolver"},
			flags: &Flags{Resolver: "some-resolver", Sudo: true, Color: true},
		},
		{
			args:  []string{"-i", "1.2.3"},
			flags: &Flags{ChefVersion: "1.2.3", Sudo: true, Color: true},
		},
		{
			args:  []string{"--install-chef", "1.2.3"},
			flags: &Flags{ChefVersion: "1.2.3", Sudo: true, Color: true},
		},
		{
			args:  []string{"-F", "some-format"},
			flags: &Flags{Format: "some-format", Sudo: true, Color: true},
		},
		{
			args:  []string{"--format", "some-format"},
			flags: &Flags{Format: "some-format", Sudo: true, Color: true},
		},
		{
			args:  []string{"-l", "some-level"},
			flags: &Flags{LogLevel: "some-level", Sudo: true, Color: true},
		},
		{
			args:  []string{"--log_level", "some-level"},
			flags: &Flags{LogLevel: "some-level", Sudo: true, Color: true},
		},
		{
			args:  []string{"-j", "some-file"},
			flags: &Flags{JSONFile: "some-file", Sudo: true, Color: true},
		},
		{
			args:  []string{"--json-attributes", "some-file"},
			flags: &Flags{JSONFile: "some-file", Sudo: true, Color: true},
		},
		{
			args:  []string{"some-recipe", "another-recipe"},
			flags: &Flags{Recipes: []string{"some-recipe", "another-recipe"}, Sudo: true, Color: true},
		},
		{
			args:  []string{"--sudo=false"},
			flags: &Flags{Sudo: false, Color: true},
		},
		{
			args:  []string{"--color=false"},
			flags: &Flags{Sudo: true, Color: false},
		},
		{
			args:  []string{"--version"},
			flags: &Flags{ShowVersion: true, Sudo: true, Color: true},
		},
		{
			args: []string{"--machine", "some-machine", "-l", "some-level", "-i", "true", "some-recipe"},
			flags: &Flags{
				Machine:     "some-machine",
				ChefVersion: "true",
				LogLevel:    "some-level",
				Recipes:     []string{"some-recipe"},
				Sudo:        true,
				Color:       true,
			},
		},
		// Check for errors
		{
			args:      []string{"-H", "some-host", "-M", "some-machine"},
			flags:     nil,
			errString: "-H, -M, -K, and -L cannot be used together",
		},
		{
			args:      []string{"-K", "some-instance", "-L"},
			flags:     nil,
			errString: "-H, -M, -K, and -L cannot be used together",
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
