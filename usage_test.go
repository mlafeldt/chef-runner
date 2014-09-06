package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	tests := []struct {
		args  []string
		flags Flags
	}{
		{
			args:  []string{},
			flags: Flags{},
		},
		{
			args:  []string{"--version"},
			flags: Flags{ShowVersion: true},
		},
		{
			args:  []string{"-H", "some-host"},
			flags: Flags{Host: "some-host"},
		},
		{
			args:  []string{"--host", "some-host"},
			flags: Flags{Host: "some-host"},
		},
		{
			args:  []string{"-M", "some-machine"},
			flags: Flags{Machine: "some-machine"},
		},
		{
			args:  []string{"--machine", "some-machine"},
			flags: Flags{Machine: "some-machine"},
		},
		{
			args:  []string{"-F", "some-format"},
			flags: Flags{Format: "some-format"},
		},
		{
			args:  []string{"--format", "some-format"},
			flags: Flags{Format: "some-format"},
		},
		{
			args:  []string{"-l", "some-level"},
			flags: Flags{LogLevel: "some-level"},
		},
		{
			args:  []string{"--log_level", "some-level"},
			flags: Flags{LogLevel: "some-level"},
		},
		{
			args:  []string{"-j", "some-file"},
			flags: Flags{JSONFile: "some-file"},
		},
		{
			args:  []string{"--json-attributes", "some-file"},
			flags: Flags{JSONFile: "some-file"},
		},
		{
			args:  []string{"some-recipe", "another-recipe"},
			flags: Flags{Recipes: []string{"some-recipe", "another-recipe"}},
		},
		{
			args: []string{"--machine", "some-machine", "-l", "some-level", "some-recipe"},
			flags: Flags{
				Machine:  "some-machine",
				LogLevel: "some-level",
				Recipes:  []string{"some-recipe"},
			},
		},
	}
	for _, test := range tests {
		flags, err := ParseFlags(test.args)
		if assert.NoError(t, err) {
			assert.Equal(t, &test.flags, flags)
		}
	}
}

func TestParseFlags_Error(t *testing.T) {
	_, err := ParseFlags([]string{"-H", "some-host", "-M", "some-machine"})
	assert.EqualError(t, err, "-H and -M cannot be used together")
}
