// Package chefsolo implements a provisioner using Chef Solo.
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
	// DefaultFormat is the default output format of Chef.
	DefaultFormat = "null"

	// DefaultLogLevel is the default log level of Chef.
	DefaultLogLevel = "info"
)

// CookbookPath is the path to the sandbox directory where cookbooks are stored.
var CookbookPath = base.SandboxPathTo("cookbooks")

// Provisioner is a provisioner based on Chef Solo.
type Provisioner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
	UseSudo    bool
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

// CreateSandbox creates the sandbox directory. This includes preparing Chef
// configuration data and cookbooks.
func (p Provisioner) CreateSandbox() error {
	funcs := []func() error{
		base.CreateSandbox,
		p.prepareJSON,
		p.prepareSoloConfig,
		p.prepareCookbooks,
	}
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// CleanupSandbox deletes the sandbox directory.
func (p Provisioner) CleanupSandbox() error {
	return base.CleanupSandbox()
}

func (p Provisioner) sudo(cmd []string) []string {
	if !p.UseSudo {
		return cmd
	}
	sudoCmd := []string{"sudo"}
	return append(sudoCmd, cmd...)
}

// Command returns the command string which will invoke the provisioner on the
// prepared machine.
func (p Provisioner) Command() []string {
	format := p.Format
	if format == "" {
		format = DefaultFormat
	}

	logLevel := p.LogLevel
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	cmd := []string{
		"chef-solo",
		"--config", base.RootPathTo("solo.rb"),
		"--json-attributes", base.RootPathTo("dna.json"),
		"--override-runlist", strings.Join(p.RunList, ","),
		"--format", format,
		"--log_level", logLevel,
	}
	return p.sudo(cmd)
}
