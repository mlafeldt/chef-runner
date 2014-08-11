// Package chefsolo implements the provisioner.Provisioner interface using Chef
// Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
	base "github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/resolver"
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

func (p Provisioner) prepareCookbooks() error {
	log.Debug("Preparing cookbooks")
	return resolver.AutoResolve(CookbookPath)
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
