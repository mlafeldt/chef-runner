// Package provisioner defines the interface that all provisioners need to
// implement. It also provides common functions shared by all provisioners. A
// provisioner is responsible for provisioning a machine with Chef.
package provisioner

import (
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/log"
)

type Provisioner interface {
	CreateSandbox() error
	CleanupSandbox() error
	Command() []string
}

var (
	SandboxPath = ".chef-runner/sandbox"
	RootPath    = "/tmp/chef-runner"
)

func SandboxPathTo(f string) string {
	return path.Join(SandboxPath, f)
}

func RootPathTo(f string) string {
	return path.Join(RootPath, f)
}

func CreateSandbox() error {
	log.Info("Preparing local files")
	log.Debug("Creating local sandbox in", SandboxPath)
	return os.MkdirAll(SandboxPath, 0755)
}

func CleanupSandbox() error {
	log.Debug("Cleaning up local sandbox in", SandboxPath)
	return os.RemoveAll(SandboxPath)
}
