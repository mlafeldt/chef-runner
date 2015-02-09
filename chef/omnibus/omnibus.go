// Package omnibus allows to install Chef using Omnibus Installer.
// See https://docs.getchef.com/install_omnibus.html
package omnibus

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/

import (
	"io/ioutil"
	"path"

	"github.com/mlafeldt/chef-runner/log"
)

// An Installer allows to install Chef.
type Installer struct {
	ChefVersion string
	SandboxPath string
	RootPath    string
	Sudo        bool
}

func (i Installer) skip() bool {
	return i.ChefVersion == "" || i.ChefVersion == "false"
}

func (i Installer) writeWrapperScript() error {
	script := path.Join(i.SandboxPath, "install-wrapper.sh")
	log.Debugf("Writing install wrapper script to %s\n", script)
	data, err := Asset("assets/install-wrapper.sh")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(script, []byte(data), 0644)
}

func (i Installer) writeOmnibusScript() error {
	script := path.Join(i.SandboxPath, "install.sh")
	log.Debugf("Writing Omnibus script to %s\n", script)
	data, err := Asset("assets/install.sh")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(script, []byte(data), 0644)
}

// PrepareFiles sets up the scripts required to install Chef.
func (i Installer) PrepareFiles() error {
	if i.skip() {
		log.Debug("Skipping setup of install scripts")
		return nil
	}
	log.Debug("Preparing install scripts")
	if err := i.writeWrapperScript(); err != nil {
		return err
	}
	return i.writeOmnibusScript()
}

// Command returns the command string to conditionally install Chef onto a
// machine.
func (i Installer) Command() []string {
	if i.skip() {
		return []string{}
	}
	cmd := []string{
		"sh",
		path.Join(i.RootPath, "install-wrapper.sh"),
		path.Join(i.RootPath, "install.sh"),
		i.ChefVersion,
	}
	if !i.Sudo {
		return cmd
	}
	return append([]string{"sudo"}, cmd...)
}
