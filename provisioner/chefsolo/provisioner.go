// Package chefsolo implements the provisioner.Provisioner interface using Chef
// Solo.
package chefsolo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlafeldt/chef-runner/cookbook"
	"github.com/mlafeldt/chef-runner/log"
	base "github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/resolver/berkshelf"
	"github.com/mlafeldt/chef-runner/resolver/librarian"
	"github.com/mlafeldt/chef-runner/resolver/rsync"
	"github.com/mlafeldt/chef-runner/util"
)

const (
	DefaultFormat   = "null"
	DefaultLogLevel = "info"
)

var (
	CookbookPath = base.SandboxPathTo("cookbooks")
)

type Provisioner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
}

func (p Provisioner) prepareJSON() error {
	log.Debug("Preparing JSON data")
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(base.SandboxPathTo("dna.json"), []byte(data), 0644)
}

func (p Provisioner) prepareSoloConfig() error {
	log.Debug("Preparing Chef Solo config")
	data := fmt.Sprintf("cookbook_path \"%s\"\n", base.RootPathTo("cookbooks"))
	data += "ssl_verify_mode :verify_peer\n"
	return ioutil.WriteFile(base.SandboxPathTo("solo.rb"), []byte(data), 0644)
}

func (p Provisioner) resolveWithBerkshelf() error {
	log.Info("Installing cookbooks with Berkshelf")
	return berkshelf.InstallCookbooks(CookbookPath)
}

func (p Provisioner) resolveWithLibrarian() error {
	log.Info("Installing cookbooks with Librarian-Chef")
	return librarian.InstallCookbooks(CookbookPath)
}

func (p Provisioner) resolveWithRsync() error {
	log.Info("Installing cookbook in current directory with rsync")
	return rsync.InstallCookbook(CookbookPath, ".")
}

func (p Provisioner) prepareCookbooks() error {
	cb, _ := cookbook.NewCookbook(".")

	// If the current folder is a cookbook and its dependencies have
	// already been resolved, only update this cookbook with rsync.
	// TODO: improve this check by comparing timestamps etc.
	if cb.Name != "" && util.FileExist(CookbookPath) {
		return p.resolveWithRsync()
	}

	if util.FileExist("Berksfile") {
		return p.resolveWithBerkshelf()
	}

	if util.FileExist("Cheffile") {
		return p.resolveWithLibrarian()
	}

	if cb.Name != "" {
		return p.resolveWithRsync()
	}

	log.Error("Berksfile, Cheffile, or metadata.rb must exist in current directory")
	return errors.New("cookbooks could not be found")
}

func (p Provisioner) CreateSandbox() error {
	if err := base.CreateSandbox(); err != nil {
		return err
	}
	if err := p.prepareJSON(); err != nil {
		return err
	}
	if err := p.prepareSoloConfig(); err != nil {
		return err
	}
	return p.prepareCookbooks()
}

func (p Provisioner) CleanupSandbox() error {
	return base.CleanupSandbox()
}

func (p Provisioner) Command() []string {
	format := p.Format
	if format == "" {
		format = DefaultFormat
	}
	logLevel := p.LogLevel
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}
	return []string{
		"sudo", "chef-solo",
		"--config", base.RootPathTo("solo.rb"),
		"--json-attributes", base.RootPathTo("dna.json"),
		"--override-runlist", strings.Join(p.RunList, ","),
		"--format", format,
		"--log_level", logLevel,
	}
}
