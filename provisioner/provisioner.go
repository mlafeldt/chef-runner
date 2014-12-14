// Package provisioner defines the interface that all provisioners need to
// implement. It also provides common functions shared by all provisioners.
package provisioner

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/log"
)

const (
	// SandboxVersion is a number that is stored inside every sandbox. This
	// number can be increased to avoid compatibility issues in case the
	// sandbox structure changes.
	SandboxVersion = "1"

	// SandboxPath is the path to the local sandbox directory where
	// chef-runner stores files that will be uploaded to a machine.
	SandboxPath = ".chef-runner/sandbox"

	// RootPath is the path on the machine where files from SandboxPath
	// will be uploaded to.
	RootPath = "/tmp/chef-runner"
)

// A Provisioner is responsible for provisioning a machine with Chef.
type Provisioner interface {
	CreateSandbox() error
	CleanupSandbox() error
	InstallCommand() []string
	ProvisionCommand() []string
}

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

	version, _ := ioutil.ReadFile(SandboxPathTo("version"))
	if string(version) != SandboxVersion {
		log.Debugf("Wiping old sandbox with version: %s\n", string(version))
		if err := os.RemoveAll(SandboxPath); err != nil {
			return err
		}
	}

	log.Debug("Creating local sandbox in", SandboxPath)
	if err := os.MkdirAll(SandboxPath, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(SandboxPathTo("version"), []byte(SandboxVersion), 0644)
}

// CleanupSandbox deletes the sandbox directory.
func CleanupSandbox() error {
	log.Debug("Cleaning up local sandbox in", SandboxPath)
	return os.RemoveAll(SandboxPath)
}
