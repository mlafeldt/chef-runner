package berkshelf

import (
	"os"

	"github.com/mlafeldt/chef-runner.go/exec"
)

func fileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func Install(path string) error {
	var cmd []string
	if fileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "berks", "install", "--path", path)
	return exec.RunCommand(cmd)
}
