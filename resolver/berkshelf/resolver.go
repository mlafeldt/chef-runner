// Package berkshelf implements a cookbook dependency resolver based on
// Berkshelf.
package berkshelf

import (
	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/util"
)

// Resolver is a cookbook dependency resolver based on Berkshelf.
type Resolver struct{}

// Command returns the command that will be executed by Resolve.
func Command(dst string) []string {
	var cmd []string
	if util.FileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "berks", "install", "--path", dst)
	return cmd
}

// Resolve runs Berkshelf to install cookbook dependencies to dst.
func (r Resolver) Resolve(dst string) error {
	return exec.RunCommand(Command(dst))
}

// String returns the resolver's name.
func (r Resolver) String() string {
	return "Berkshelf resolver"
}
