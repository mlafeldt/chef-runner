// Package chefsolo implements a provisioner using Chef Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/chef/omnibus"
	"github.com/mlafeldt/chef-runner/log"
)

const (
	// DefaultFormat is the default output format of Chef.
	DefaultFormat = "doc"

	// DefaultLogLevel is the default log level of Chef.
	DefaultLogLevel = "info"
)

// Provisioner is a provisioner based on Chef Solo.
type Provisioner struct {
	RunList     []string
	Attributes  string
	Format      string
	LogLevel    string
	UseSudo     bool
	ChefVersion string

	SandboxPath string
	RootPath    string
}

func (p Provisioner) prepareSandbox() error {
	log.Info("Preparing local files")
	log.Debug("Creating local sandbox in", p.SandboxPath)
	return os.MkdirAll(p.SandboxPath, 0755)
}

func (p Provisioner) prepareJSON() error {
	log.Debug("Preparing JSON data")
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(path.Join(p.SandboxPath, "dna.json"), []byte(data), 0644)
}

func (p Provisioner) prepareSoloConfig() error {
	log.Debug("Preparing Chef Solo config")
	data := fmt.Sprintf("cookbook_path \"%s\"\n", path.Join(p.RootPath, "cookbooks"))
	data += "ssl_verify_mode :verify_peer\n"
	return ioutil.WriteFile(path.Join(p.SandboxPath, "solo.rb"), []byte(data), 0644)
}

func (p Provisioner) prepareInstallScripts() error {
	i := omnibus.Installer{
		ChefVersion: p.ChefVersion,
		ScriptPath:  p.SandboxPath,
	}
	return i.PrepareScripts()
}

// PrepareFiles creates the sandbox directory. This includes preparing Chef
// configuration data and cookbooks.
func (p Provisioner) PrepareFiles() error {
	funcs := []func() error{
		p.prepareSandbox,
		p.prepareJSON,
		p.prepareSoloConfig,
		p.prepareInstallScripts,
	}
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// InstallCommand returns the command string to conditionally install Chef onto
// a machine.
func (p Provisioner) InstallCommand() []string {
	i := omnibus.Installer{
		ChefVersion: p.ChefVersion,
		ScriptPath:  p.RootPath,
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
		"--config", path.Join(p.RootPath, "solo.rb"),
		"--json-attributes", path.Join(p.RootPath, "dna.json"),
		"--format", format,
		"--log_level", logLevel,
	}

	if len(p.RunList) > 0 {
		cmd = append(cmd, "--override-runlist", strings.Join(p.RunList, ","))
	}

	return p.sudo(cmd)
}
