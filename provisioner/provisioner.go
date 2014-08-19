// Package provisioner defines the interface that all provisioners need to
// implement. It also provides common functions shared by all provisioners.
package provisioner

import (
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/log"
)

// A Provisioner is responsible for provisioning a machine with Chef.
type Provisioner interface {
	CreateSandbox() error
	CleanupSandbox() error
	Command() []string
}

var (
	// SandboxPath is the path to the local sandbox directory where
	// chef-runner stores files that will be uploaded to a machine.
	SandboxPath = ".chef-runner/sandbox"

	// RootPath is the path on the machine where files from SandboxPath
	// will be uploaded to.
	RootPath = "/tmp/chef-runner"
)

// SandboxPathTo returns a path relative to SandboxPath.
func SandboxPathTo(elem ...string) string {
	slice := append([]string{SandboxPath}, elem...)
	return path.Join(slice...)
}

// RootPathTo returns a path relative to RootPath.
func RootPathTo(elem ...string) string {
	slice := append([]string{RootPath}, elem...)
	return path.Join(slice...)
}

// CreateSandbox creates the sandbox directory.
func CreateSandbox() error {
	log.Info("Preparing local files")
	log.Debug("Creating local sandbox in", SandboxPath)
	return os.MkdirAll(SandboxPath, 0755)
}

// CleanupSandbox deletes the sandbox directory.
func CleanupSandbox() error {
	log.Debug("Cleaning up local sandbox in", SandboxPath)
	return os.RemoveAll(SandboxPath)
}
