// Package omnibus allows to install Chef using Omnibus Installer.
// See https://docs.getchef.com/install_omnibus.html
package omnibus

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/util"
)

// ScriptURL is the URL of the Omnibus install script.
var ScriptURL = "https://www.opscode.com/chef/install.sh"

// An Installer allows to install Chef.
type Installer struct {
	ChefVersion string
	ScriptPath  string
}

func (i Installer) skip() bool {
	return i.ChefVersion == "" || i.ChefVersion == "false"
}

func (i Installer) scriptPathTo(elem ...string) string {
	slice := append([]string{i.ScriptPath}, elem...)
	return path.Join(slice...)
}

func (i Installer) writeWrapperScript() error {
	script := i.scriptPathTo("install-wrapper.sh")
	log.Debugf("Writing install wrapper script to %s\n", script)
	return ioutil.WriteFile(script, []byte(wrapperScript), 0644)
}

func (i Installer) downloadOmnibusScript() error {
	script := i.scriptPathTo("install.sh")
	if util.FileExist(script) {
		log.Debugf("Omnibus script already downloaded to %s\n", script)
		return nil
	}
	log.Debugf("Downloading Omnibus script from %s to %s\n", ScriptURL, script)
	return util.DownloadFile(script, ScriptURL)
}

// PrepareScripts sets up the scripts required to install Chef.
func (i Installer) PrepareScripts() error {
	if i.skip() {
		log.Debug("Skipping setup of install scripts")
		return nil
	}
	log.Debug("Preparing install scripts")
	if err := os.MkdirAll(i.ScriptPath, 0755); err != nil {
		return err
	}
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
		i.scriptPathTo("install-wrapper.sh"),
		i.scriptPathTo("install.sh"),
		i.ChefVersion,
	}
}
