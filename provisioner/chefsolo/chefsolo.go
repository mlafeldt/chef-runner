package chefsolo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner.go/berkshelf"
	"github.com/mlafeldt/chef-runner.go/cookbook"
	"github.com/mlafeldt/chef-runner.go/rsync"
	"github.com/mlafeldt/chef-runner.go/util"
)

const (
	SandboxPath = ".chef-runner"

	// TODO: change prefix from /vagrant to /tmp and explicitly copy files
	// there in order to get rid of the Vagrant dependency
	RootPath = "/vagrant/" + SandboxPath

	DefaultFormat   = "null"
	DefaultLogLevel = "info"
)

func sandboxPath(f string) string {
	return path.Join(SandboxPath, f)
}

func rootPath(f string) string {
	return path.Join(RootPath, f)
}

type Provisoner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
}

func (p *Provisoner) prepareJSON() error {
	data := "{}\n"
	if p.Attributes != "" {
		data = p.Attributes
	}
	return ioutil.WriteFile(sandboxPath("dna.json"), []byte(data), 0644)
}

func (p *Provisoner) prepareSoloConfig() error {
	data := fmt.Sprintf("cookbook_path \"%s\"\n", rootPath("cookbooks"))
	return ioutil.WriteFile(sandboxPath("solo.rb"), []byte(data), 0644)
}

func (p *Provisoner) prepareCookbooks() error {
	cookbookPath := sandboxPath("cookbooks")
	if !util.FileExist(cookbookPath) {
		return berkshelf.Install(cookbookPath)
	}
	cb, err := cookbook.New(".")
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

func (p *Provisoner) CreateSandbox() error {
	if err := os.MkdirAll(SandboxPath, 0755); err != nil {
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

func (p *Provisoner) CleanupSandbox() error {
	return os.RemoveAll(SandboxPath)
}

func (p *Provisoner) Command() []string {
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
		"--config", rootPath("solo.rb"),
		"--json-attributes", rootPath("dna.json"),
		"--override-runlist", strings.Join(p.RunList, ","),
		"--format", format,
		"--log_level", logLevel,
	}
}
