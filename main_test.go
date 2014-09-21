package main

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/stretchr/testify/assert"
)

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
		assert.Equal(t, level, logLevel())
	}
}

func TestBuildRunList(t *testing.T) {
	tests := []struct {
		cookbook  string
		recipes   []string
		runList   []string
		errString string
	}{
		{
			cookbook: "cats",
			recipes:  []string{},
			runList:  []string{},
		},
		{
			cookbook: "cats",
			recipes:  []string{"::foo"},
			runList:  []string{"cats::foo"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"recipes/foo.rb"},
			runList:  []string{"cats::foo"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"./recipes//foo.rb"},
			runList:  []string{"cats::foo"},
		},
		{
			cookbook: "",
			recipes:  []string{"dogs::bar"},
			runList:  []string{"dogs::bar"},
		},
		{
			cookbook: "",
			recipes:  []string{"dogs"},
			runList:  []string{"dogs"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"recipes/foo.rb", "::bar", "dogs::baz"},
			runList:  []string{"cats::foo", "cats::bar", "dogs::baz"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"recipes/foo.rb,::bar,dogs::baz"},
			runList:  []string{"cats::foo", "cats::bar", "dogs::baz"},
		},
		// Check for errors
		{
			cookbook:  "",
			recipes:   []string{"::foo"},
			runList:   nil,
			errString: "cookbook name required",
		},
		{
			cookbook:  "",
			recipes:   []string{"recipes/foo.rb"},
			runList:   nil,
			errString: "cookbook name required",
		},
	}
	for _, test := range tests {
		runList, err := buildRunList(test.recipes, test.cookbook)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.runList, runList)
	}
}
