package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/mlafeldt/chef-runner.go/exec"
	"github.com/mlafeldt/chef-runner.go/vagrant"
)

const (
	CookbookPath    = "vendor/cookbooks"
	VagrantChefPath = "/tmp/vagrant-chef-1"
)

type Options struct {
	Host     string `short:"H" long:"host" description:"Set hostname for direct SSH access" value-name:"NAME"`
	Machine  string `short:"M" long:"machine" description:"Set name of Vagrant virtual machine" value-name:"NAME"`
	Format   string `short:"F" long:"format" default:"null" description:"Set output format" value-name:"FORMAT"`
	LogLevel string `short:"l" long:"log_level" default:"info" description:"Set log level" value-name:"LEVEL"`
	JsonFile string `short:"j" long:"json-attributes" description:"Load attributes from a JSON file" value-name:"FILE"`
}

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

func parseFlags(opts *Options, argv []string) []string {
	args, err := flags.ParseArgs(opts, argv)
	if err != nil {
		// --help is not an error
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		// FIXME: show internal parser errors
		os.Exit(1)
	}

	return args
}

func fileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func rsyncCookbook(cookbookName, installDir string) error {
	// TODO: filter files more diligently
	files, err := filepath.Glob("[a-zA-Z]*")
	if err != nil {
		return err
	}
	cmd := []string{"rsync", "--archive", "--delete", "--exclude=" + installDir}
	cmd = append(cmd, files...)
	cmd = append(cmd, path.Join(installDir, cookbookName))
	return exec.RunCommand(cmd)
}

func berkshelf(installDir string) error {
	var cmd []string
	if fileExist("Gemfile") {
		cmd = []string{"bundle", "exec"}
	}
	cmd = append(cmd, "berks", "install", "--path", installDir)
	return exec.RunCommand(cmd)
}

// Install cookbook dependencies with Berkshelf. If the cookbooks are already
// in place, use lightning-fast rsync to update the current cookbook only.
func installCookbooks(cookbookName, installDir string) error {
	if fileExist(installDir) {
		rsyncCookbook(cookbookName, installDir)
	}
	return berkshelf(installDir)
}

func openSSH(host, command string) error {
	return exec.RunCommand([]string{"ssh", host, "-c", command})
}

func provision(opts Options, runlist string) error {
	config_file := VagrantChefPath + "/solo.rb"
	json_file := VagrantChefPath + "/dna.json"
	cookbooks_path := "/vagrant/" + CookbookPath

	setup_dir := fmt.Sprintf("sudo mkdir -p %s", VagrantChefPath)
	setup_config := fmt.Sprintf("test -f %s || echo 'cookbook_path \"%s\"' | sudo tee %s >/dev/null", config_file, cookbooks_path, config_file)
	setup_json := fmt.Sprintf("test -f %s || echo '{}' | sudo tee %s >/dev/null", json_file, json_file)

	if opts.JsonFile != "" {
		json_file = "/vagrant/" + opts.JsonFile
	}

	run_chef_solo := fmt.Sprintf("sudo chef-solo --config=%s --json-attributes=%s --override-runlist=%s --format=%s --log_level=%s",
		config_file, json_file, runlist, opts.Format, opts.LogLevel)

	cmd := strings.Join([]string{setup_dir, setup_config, setup_json, run_chef_solo}, " && ")
	// fmt.Println(cmd)

	var err error
	if opts.Host != "" {
		err = openSSH(opts.Host, cmd)
	} else {
		err = vagrant.RunCommand(opts.Machine, cmd)
	}
	return err
}

func cookbookName(cookbookPath string) string {
	base := path.Base(cookbookPath)
	if strings.HasPrefix(base, "chef-") {
		return strings.TrimPrefix(base, "chef-")
	}
	if strings.HasSuffix(base, "-cookbook") {
		return strings.TrimSuffix(base, "-cookbook")
	}
	return base
}

func baseName(s, suffix string) string {
	base := path.Base(s)
	if suffix != "" {
		base = strings.TrimSuffix(base, suffix)
	}
	return base
}

func buildRunList(cookbookName string, recipes []string) string {
	if len(recipes) == 0 {
		return cookbookName + "::default"
	}

	var runlist []string
	for _, r := range recipes {
		var recipeName string
		if strings.Contains(r, "::") {
			recipeName = r
		} else if path.Dir(r) == "recipes" && path.Ext(r) == ".rb" {
			recipeName = cookbookName + "::" + baseName(r, ".rb")
		} else {
			recipeName = cookbookName + "::" + r
		}
		runlist = append(runlist, recipeName)
	}
	return strings.Join(runlist, ",")
}

func main() {
	var opts Options

	log.SetFlags(0)

	args := parseFlags(&opts, os.Args[1:])
	if opts.Host != "" && opts.Machine != "" {
		log.Fatal("error: --host and --machine cannot be used together")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cookbookName := cookbookName(cwd)

	runlist := buildRunList(cookbookName, args)
	fmt.Println("Run List is", runlist)

	if err := installCookbooks(cookbookName, CookbookPath); err != nil {
		log.Fatal(err)
	}
	if err := provision(opts, runlist); err != nil {
		log.Fatal(err)
	}
}
