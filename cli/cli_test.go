package cli_test

import (
	"testing"

	. "github.com/mlafeldt/chef-runner/cli"
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
			args:  []string{"--version"},
			flags: &Flags{ShowVersion: true},
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
			args:  []string{"--ssh-option", "x=1", "--ssh-option", "y 2 3"},
			flags: &Flags{SSHOptions: []string{"x=1", "y 2 3"}},
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
			args:  []string{"some-recipe", "another-recipe"},
			flags: &Flags{Recipes: []string{"some-recipe", "another-recipe"}},
		},
		{
			args: []string{"--machine", "some-machine", "-l", "some-level", "some-recipe"},
			flags: &Flags{
				Machine:  "some-machine",
				LogLevel: "some-level",
				Recipes:  []string{"some-recipe"},
			},
		},
		// Check for errors
		{
			args:      []string{"-H", "some-host", "-M", "some-machine"},
			flags:     nil,
			errString: "-H and -M cannot be used together",
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
