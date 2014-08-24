// Package bundler helps to run external commands with Bundler if the
// environment indicates that Bundler should be used.
package bundler

import (
	"os/exec"

	"github.com/mlafeldt/chef-runner/util"
)

func useBundler() bool {
	if _, err := exec.LookPath("bundle"); err != nil {
		// Bundler not installed
		return false
	}
	if !util.FileExist("Gemfile") {
		// No Gemfile found
		return false
	}
	return true
}

// Command prepends `bundle exec` to the passed command if the environment
// indicates that Bundler should be used.
func Command(args []string) []string {
	if !useBundler() {
		return args
	}
	return append([]string{"bundle", "exec"}, args...)
}
