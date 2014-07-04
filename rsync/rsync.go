package rsync

import (
	"github.com/mlafeldt/chef-runner.go/exec"
)

type Options struct {
	Archive bool
	Delete  bool
	Verbose bool
	Exclude []string
}

func Copy(src []string, dst string, opts Options) error {
	cmd := []string{"rsync"}
	if opts.Archive {
		cmd = append(cmd, "--archive")
	}
	if opts.Delete {
		cmd = append(cmd, "--delete")
	}
	if opts.Verbose {
		cmd = append(cmd, "--verbose")
	}
	for _, x := range opts.Exclude {
		cmd = append(cmd, "--exclude", x)
	}
	cmd = append(cmd, src...)
	cmd = append(cmd, dst)
	return exec.RunCommand(cmd)
}
