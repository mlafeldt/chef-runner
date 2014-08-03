package librarian

import (
	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/util"
)

func InstallCommand(path string) []string {
	var cmd []string
	if util.FileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "librarian-chef", "install", "--path", path)
	return cmd
}

func Install(path string) error {
	return exec.RunCommand(InstallCommand(path))
}
