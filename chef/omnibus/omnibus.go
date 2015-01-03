// Package omnibus allows to install Chef using Omnibus Installer.
// See https://docs.getchef.com/install_omnibus.html
package omnibus

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/

import (
	"io/ioutil"
	"path"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/util"
)

// ScriptURL is the URL of the Omnibus install script.
var ScriptURL = "https://www.opscode.com/chef/install.sh"

// An Installer allows to install Chef.
type Installer struct {
	ChefVersion string
	SandboxPath string
	RootPath    string
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

func (i Installer) downloadOmnibusScript() error {
	script := path.Join(i.SandboxPath, "install.sh")
	if util.FileExist(script) {
		log.Debugf("Omnibus script already downloaded to %s\n", script)
		return nil
	}
	log.Debugf("Downloading Omnibus script from %s to %s\n", ScriptURL, script)
	return util.DownloadFile(script, ScriptURL)
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
	return i.downloadOmnibusScript()
}

// Command returns the command string to conditionally install Chef onto a
// machine.
func (i Installer) Command() []string {
	if i.skip() {
		return []string{}
	}
	return []string{
		"sudo",
		"sh",
		path.Join(i.RootPath, "install-wrapper.sh"),
		path.Join(i.RootPath, "install.sh"),
		i.ChefVersion,
	}
}
