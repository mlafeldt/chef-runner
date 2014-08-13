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
		cookbookName string
		recipes      []string
		runList      []string
		errString    string
	}{
		{
			cookbookName: "cats",
			recipes:      []string{},
			runList:      []string{"cats::default"},
		},
		{
			cookbookName: "cats",
			recipes:      []string{"recipes/foo.rb"},
			runList:      []string{"cats::foo"},
		},
		{
			cookbookName: "cats",
			recipes:      []string{"./recipes//foo.rb"},
			runList:      []string{"cats::foo"},
		},
		{
			cookbookName: "cats",
			recipes:      []string{"foo"},
			runList:      []string{"cats::foo"},
		},
		{
			cookbookName: "",
			recipes:      []string{"dogs::bar"},
			runList:      []string{"dogs::bar"},
		},
		{
			cookbookName: "cats",
			recipes:      []string{"recipes/foo.rb", "bar", "dogs::baz"},
			runList:      []string{"cats::foo", "cats::bar", "dogs::baz"},
		},
		// Check for errors
		{
			cookbookName: "",
			recipes:      []string{},
			runList:      nil,
			errString:    "cookbook name required",
		},
		{
			cookbookName: "",
			recipes:      []string{"foo"},
			runList:      nil,
			errString:    "cookbook name required",
		},
		{
			cookbookName: "",
			recipes:      []string{"recipes/foo.rb"},
			runList:      nil,
			errString:    "cookbook name required",
		},
	}
	for _, test := range tests {
		runList, err := buildRunList(test.cookbookName, test.recipes)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.runList, runList)
	}
}

func TestVersionString(t *testing.T) {
	GitVersion = ""
	assert.Equal(t, Version, VersionString())

	GitVersion = "some-git-version"
	assert.Equal(t, GitVersion, VersionString())
}
