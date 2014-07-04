package vagrant

import (
	"os"

	"github.com/mlafeldt/chef-runner.go/exec"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

func RunCommand(machine, command string) error {
	if machine == "" {
		machine = "default"
	}
	return exec.RunCommand([]string{"vagrant", "ssh", machine, "-c", command})
}
