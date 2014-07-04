package berkshelf

import (
	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/util"
)

func Install(path string) error {
	var cmd []string
	if util.FileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "berks", "install", "--path", path)
	return exec.RunCommand(cmd)
}
