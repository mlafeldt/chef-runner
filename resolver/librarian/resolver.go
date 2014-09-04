// Package librarian implements a cookbook dependency resolver based on
// Librarian-Chef.
package librarian

import (
	"github.com/mlafeldt/chef-runner/bundler"
	"github.com/mlafeldt/chef-runner/exec"
)

// Resolver is a cookbook dependency resolver based on Librarian-Chef.
type Resolver struct{}

// Command returns the command that will be executed by Resolve.
func Command(dst string) []string {
	cmd := []string{"librarian-chef", "install", "--path", dst}
	return bundler.Command(cmd)
}

// Resolve runs Librarian-Chef to install cookbook dependencies to dst.
func (r Resolver) Resolve(dst string) error {
	return exec.RunCommand(Command(dst))
}

// String returns the resolver's name.
func (r Resolver) String() string {
	return "Librarian-Chef resolver"
}
