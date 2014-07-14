package provisioner

import (
	"os"
	"path"
)

var (
	SandboxPath = ".chef-runner"

	// TODO: change prefix from /vagrant to /tmp and explicitly copy files
	// there in order to get rid of the Vagrant dependency
	RootPath = "/vagrant/" + SandboxPath
)

func SandboxPathTo(f string) string {
	return path.Join(SandboxPath, f)
}

func RootPathTo(f string) string {
	return path.Join(RootPath, f)
}

func CreateSandbox() error {
	return os.MkdirAll(SandboxPath, 0755)
}

func CleanupSandbox() error {
	return os.RemoveAll(SandboxPath)
}
