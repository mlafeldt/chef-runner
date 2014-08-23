// Package bundler helps to run external commands with Bundler if the
// environment indicates that Bundler should be used.
package bundler

import "github.com/mlafeldt/chef-runner/util"

func useBundler() bool {
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
