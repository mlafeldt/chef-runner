// Package chefsolo implements the provisioner.Provisoner interface using Chef
// Solo.
package chefsolo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlafeldt/chef-runner/cookbook/metadata"
	"github.com/mlafeldt/chef-runner/log"
	. "github.com/mlafeldt/chef-runner/provisioner"
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
	CookbookPath = SandboxPathTo("cookbooks")
)

type Provisoner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
}

func (p Provisoner) prepareJSON() error {
	log.Debug("Preparing JSON data")
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(SandboxPathTo("dna.json"), []byte(data), 0644)
}

func (p Provisoner) prepareSoloConfig() error {
	log.Debug("Preparing Chef Solo config")
	data := fmt.Sprintf("cookbook_path \"%s\"\n", RootPathTo("cookbooks"))
	return ioutil.WriteFile(SandboxPathTo("solo.rb"), []byte(data), 0644)
}

func (p Provisoner) resolveWithBerkshelf() error {
	log.Info("Installing cookbooks with Berkshelf")
	return berkshelf.InstallCookbooks(CookbookPath)
}

func (p Provisoner) resolveWithLibrarian() error {
	log.Info("Installing cookbooks with Librarian-Chef")
	return librarian.InstallCookbooks(CookbookPath)
}

func (p Provisoner) resolveWithRsync() error {
	log.Info("Installing cookbook in current directory with rsync")
	return rsync.InstallCookbook(CookbookPath, ".")
}

func (p Provisoner) prepareCookbooks() error {
	// If the current folder is a cookbook and its dependencies have
	// already been resolved, only update this cookbook with rsync.
	// TODO: improve this check by comparing timestamps etc.
	if util.FileExist(metadata.Filename) && util.FileExist(CookbookPath) {
		return p.resolveWithRsync()
	}

	if util.FileExist("Berksfile") {
		return p.resolveWithBerkshelf()
	}

	if util.FileExist("Cheffile") {
		return p.resolveWithLibrarian()
	}

	if util.FileExist(metadata.Filename) {
		return p.resolveWithRsync()
	}

	log.Error("Berksfile, Cheffile, or metadata.rb must exist in current directory")
	return errors.New("cookbooks could not be found")
}

func (p Provisoner) CreateSandbox() error {
	if err := CreateSandbox(); err != nil {
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

func (p Provisoner) CleanupSandbox() error {
	return CleanupSandbox()
}

func (p Provisoner) Command() []string {
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
		"--config", RootPathTo("solo.rb"),
		"--json-attributes", RootPathTo("dna.json"),
		"--override-runlist", strings.Join(p.RunList, ","),
		"--format", format,
		"--log_level", logLevel,
	}
}
