package vagrant

import (
	"github.com/mlafeldt/chef-runner.go/exec"
)

func RunCommand(machine, command string) error {
	if machine == "" {
		machine = "default"
	}
	return exec.RunCommand([]string{"vagrant", "ssh", machine, "-c", command})
}
