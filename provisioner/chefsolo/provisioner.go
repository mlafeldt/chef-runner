// Package chefsolo implements a provisioner using Chef Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
	base "github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/resolver"
	"github.com/mlafeldt/chef-runner/util"
)

const (
	// DefaultFormat is the default output format of Chef.
	DefaultFormat = "doc"

	// DefaultLogLevel is the default log level of Chef.
	DefaultLogLevel = "info"

	// OmnibusScriptURL is the URL of the Omnibus install script.
	OmnibusScriptURL = "https://www.opscode.com/chef/install.sh"
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

func (p Provisioner) downloadOmnibusScript() error {
	if len(p.InstallCommand()) == 0 {
		log.Debug("Skipping download of Omnibus script")
		return nil
	}

	script := base.SandboxPathTo("install.sh")
	if util.FileExist(script) {
		log.Debugf("Omnibus script already downloaded to %s\n", script)
		return nil
	}

	log.Debugf("Downloading Omnibus script to %s\n", script)
	return util.DownloadFile(script, OmnibusScriptURL)
}

// CreateSandbox creates the sandbox directory. This includes preparing Chef
// configuration data and cookbooks.
func (p Provisioner) CreateSandbox() error {
	funcs := []func() error{
		base.CreateSandbox,
		p.prepareJSON,
		p.prepareSoloConfig,
		p.prepareCookbooks,
		p.downloadOmnibusScript,
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

// InstallCommand returns the command string to conditionally install a Chef
// Omnibus package onto a machine.
func (p Provisioner) InstallCommand() []string {
	versionFile := "/opt/chef/version-manifest.txt"
	installCmd := fmt.Sprintf(`sudo sh %s`, base.RootPathTo("install.sh"))

	switch p.ChefVersion {
	case "", "false":
		// Do nothing
		return []string{}
	case "latest":
		// Always install latest version of Chef
		return []string{installCmd}
	case "true":
		// Only install Chef if not already installed
		checkCmd := fmt.Sprintf(`test -f %s ||`, versionFile)
		cmd := strings.Join([]string{checkCmd, installCmd}, " ")
		return []string{cmd}
	default:
		// Install specific Chef version if that version is not already installed
		checkCmd := fmt.Sprintf(`test "$(head -n1 %s 2>/dev/null | cut -d" " -f2)" = "%s" ||`,
			versionFile, p.ChefVersion)
		cmd := strings.Join([]string{checkCmd, installCmd, "-v", p.ChefVersion}, " ")
		return []string{cmd}
	}
}

func (p Provisioner) sudo(args []string) []string {
	if !p.UseSudo {
		return args
	}
	return append([]string{"sudo"}, args...)
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
