// Package berkshelf implements a cookbook dependency resolver based on
// Berkshelf.
package berkshelf

import (
	"fmt"
	"os"
	"strings"

	"github.com/mlafeldt/chef-runner/bundler"
	"github.com/mlafeldt/chef-runner/exec"
)

// Resolver is a cookbook dependency resolver based on Berkshelf.
type Resolver struct{}

// Command returns the command that will be executed by Resolve.
func Command(dst string) []string {
	code := []string{
		`require "berkshelf";`,
		`b = Berkshelf::Berksfile.from_file("Berksfile");`,
		`Berkshelf::Berksfile.method_defined?(:vendor)`, `?`,
		fmt.Sprintf(`b.vendor("%s")`, dst), `:`,
		fmt.Sprintf(`b.install(:path => "%s")`, dst),
	}
	cmd := append([]string{"ruby", "-e"}, strings.Join(code, " "))
	return bundler.Command(cmd)
}

// Resolve runs Berkshelf to install cookbook dependencies to dst.
func (r Resolver) Resolve(dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	return exec.RunCommand(Command(dst))
}

// Name returns the resolver's name.
func (r Resolver) Name() string {
	return "Berkshelf"
}
