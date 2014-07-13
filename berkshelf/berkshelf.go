package berkshelf

import (
	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/util"
)

func InstallCommand(path string) []string {
	var cmd []string
	if util.FileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "berks", "install", "--path", path)
	return cmd
}

func Install(path string) error {
	return exec.RunCommand(InstallCommand(path))
}
