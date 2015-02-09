// Package chefsolo implements a provisioner using Chef Solo.
package chefsolo

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

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
	SandboxPath string
	RootPath    string
	Sudo        bool
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

// PrepareFiles prepares Chef configuration data.
func (p Provisioner) PrepareFiles() error {
	if err := p.prepareJSON(); err != nil {
		return err
	}
	return p.prepareSoloConfig()
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
		"--config", path.Join(p.RootPath, "solo.rb"),
		"--json-attributes", path.Join(p.RootPath, "dna.json"),
		"--format", format,
		"--log_level", logLevel,
	}

	if len(p.RunList) > 0 {
		cmd = append(cmd, "--override-runlist", strings.Join(p.RunList, ","))
	}

	if !p.Sudo {
		return cmd
	}
	return append([]string{"sudo"}, cmd...)
}
