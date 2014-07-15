package chefsolo

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner.go/berkshelf"
	"github.com/mlafeldt/chef-runner.go/cookbook"
	. "github.com/mlafeldt/chef-runner.go/provisioner"
	"github.com/mlafeldt/chef-runner.go/rsync"
	"github.com/mlafeldt/chef-runner.go/util"
)

const (
	DefaultFormat   = "null"
	DefaultLogLevel = "info"
)

type Provisoner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
}

func (p Provisoner) prepareJSON() error {
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(SandboxPathTo("dna.json"), []byte(data), 0644)
}

func (p Provisoner) prepareSoloConfig() error {
	data := fmt.Sprintf("cookbook_path \"%s\"\n", RootPathTo("cookbooks"))
	return ioutil.WriteFile(SandboxPathTo("solo.rb"), []byte(data), 0644)
}

func (p Provisoner) prepareCookbooks() error {
	cookbookPath := SandboxPathTo("cookbooks")
	if !util.FileExist(cookbookPath) {
		return berkshelf.Install(cookbookPath)
	}
	cb, err := cookbook.NewCookbook(".")
	if err != nil {
		return err
	}
	files, err := cb.Files()
	if err != nil {
		return err
	}
	c := rsync.Client{Archive: true, Delete: true, Verbose: true}
	return c.Copy(files, path.Join(cookbookPath, cb.Name))
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
