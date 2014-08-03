package librarian

import (
	"os"
	"path"
	"path/filepath"

	"github.com/mlafeldt/chef-runner/exec"
	"github.com/mlafeldt/chef-runner/util"
)

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

func InstallCookbooks(dst string) error {
	if err := exec.RunCommand(Command(dst)); err != nil {
		return err
	}
	return removeTempFiles(dst)
}
