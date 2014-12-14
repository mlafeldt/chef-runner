// Package chefsolo implements a provisioner using Chef Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlafeldt/chef-runner/chef/omnibus"
	"github.com/mlafeldt/chef-runner/log"
	base "github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/resolver"
)

const (
	// DefaultFormat is the default output format of Chef.
	DefaultFormat = "doc"

	// DefaultLogLevel is the default log level of Chef.
	DefaultLogLevel = "info"
)

// CookbookPath is the path to the sandbox directory where cookbooks are stored.
var CookbookPath = base.SandboxPathTo("cookbooks")

// Provisioner is a provisioner based on Chef Solo.
type Provisioner struct {
	RunList     []string
	Attributes  string
	Format      string
	LogLevel    string
	UseSudo     bool
	ChefVersion string
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

func (p Provisioner) prepareInstallScripts() error {
	i := omnibus.Installer{
		ChefVersion: p.ChefVersion,
		ScriptPath:  base.SandboxPathTo("scripts"),
	}
	return i.PrepareScripts()
}

// CreateSandbox creates the sandbox directory. This includes preparing Chef
// configuration data and cookbooks.
func (p Provisioner) CreateSandbox() error {
	funcs := []func() error{
		base.CreateSandbox,
		p.prepareJSON,
		p.prepareSoloConfig,
		p.prepareCookbooks,
		p.prepareInstallScripts,
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

// InstallCommand returns the command string to conditionally install Chef onto
// a machine.
func (p Provisioner) InstallCommand() []string {
	i := omnibus.Installer{
		ChefVersion: p.ChefVersion,
		ScriptPath:  base.RootPathTo("scripts"),
	}
	return i.Command()
}

func (p Provisioner) sudo(args []string) []string {
	if !p.UseSudo {
		return args
	}
	return append([]string{"sudo"}, args...)
}

// ProvisionCommand returns the command string which will invoke the
// provisioner on the prepared machine.
func (p Provisioner) ProvisionCommand() []string {
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
		"--format", format,
		"--log_level", logLevel,
	}

	if len(p.RunList) > 0 {
		cmd = append(cmd, "--override-runlist", strings.Join(p.RunList, ","))
	}

	return p.sudo(cmd)
}
