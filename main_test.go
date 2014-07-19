package main

import (
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	tests := map[string]int{
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
	}{
		{"cats", []string{}, []string{"cats::default"}},
		{"cats", []string{"recipes/foo.rb"}, []string{"cats::foo"}},
		{"cats", []string{"./recipes//foo.rb"}, []string{"cats::foo"}},
		{"cats", []string{"foo"}, []string{"cats::foo"}},
		{"cats", []string{"dogs::bar"}, []string{"dogs::bar"}},
		{"cats", []string{"recipes/foo.rb", "bar", "dogs::baz"},
			[]string{"cats::foo", "cats::bar", "dogs::baz"}},
	}
	for _, test := range tests {
		runList := buildRunList(test.cookbookName, test.recipes)
		assert.Equal(t, test.runList, runList)
	}
}
