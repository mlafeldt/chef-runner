// Package librarian implements a cookbook dependency resolver based on
// Librarian-Chef.
package librarian

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/util"
)

// Resolver is a cookbook dependency resolver based on Librarian-Chef.
type Resolver struct{}

// Command returns the command that will be executed by Resolve.
func Command(dst string) []string {
	var cmd []string
	if util.FileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "librarian-chef", "install", "--path", dst)
	return cmd
}

func removeTempFiles(dst string) error {
	tmpDirs, err := filepath.Glob(path.Join(dst, "*", "tmp", "librarian"))
	if err != nil {
		return err
	}
	for _, dir := range tmpDirs {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

// Resolve runs Librarian-Chef to install cookbook dependencies to dst. It also
// removes temporary Librarian-Chef files from the installed cookbooks.
func (r Resolver) Resolve(dst string) error {
	if err := exec.RunCommand(Command(dst)); err != nil {
		return err
	}
	return removeTempFiles(dst)
}

// String returns the resolver's name.
func (r Resolver) String() string {
	return "Librarian-Chef resolver"
}
