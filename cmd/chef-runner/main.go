package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner.go/berkshelf"
	"github.com/mlafeldt/chef-runner.go/cookbook"
	"github.com/mlafeldt/chef-runner.go/openssh"
	"github.com/mlafeldt/chef-runner.go/rsync"
	"github.com/mlafeldt/chef-runner.go/util"
	"github.com/mlafeldt/chef-runner.go/vagrant"
)

const (
	CookbookPath    = "vendor/cookbooks"
	VagrantChefPath = "/tmp/vagrant-chef-1"
)

type SSHClient interface {
	RunCommand(command string) error
}

func buildRunList(cb *cookbook.Cookbook, recipes []string) string {
	if len(recipes) == 0 {
		return cb.Name + "::default"
	}

	var runlist []string
	for _, r := range recipes {
		var recipeName string
		if strings.Contains(r, "::") {
			recipeName = r
		} else if path.Dir(r) == "recipes" && path.Ext(r) == ".rb" {
			recipeName = cb.Name + "::" + util.BaseName(r, ".rb")
		} else {
			recipeName = cb.Name + "::" + r
		}
		runlist = append(runlist, recipeName)
	}
	return strings.Join(runlist, ",")
}

func installCookbooks(cb *cookbook.Cookbook, installDir string) error {
	if !util.FileExist(installDir) {
		return berkshelf.Install(installDir)
	}
	files, err := cb.Files()
	if err != nil {
		return err
	}
	c := rsync.Client{Archive: true, Delete: true, Verbose: true}
	return c.Copy(files, path.Join(installDir, cb.Name))
}

func provision(client SSHClient, format, logLevel, jsonFile string, runlist string) error {
	config_file := VagrantChefPath + "/solo.rb"
	json_file := VagrantChefPath + "/dna.json"
	cookbooks_path := "/vagrant/" + CookbookPath

	setup_dir := fmt.Sprintf("sudo mkdir -p %s", VagrantChefPath)
	setup_config := fmt.Sprintf("test -f %s || echo 'cookbook_path \"%s\"' | sudo tee %s >/dev/null", config_file, cookbooks_path, config_file)
	setup_json := fmt.Sprintf("test -f %s || echo '{}' | sudo tee %s >/dev/null", json_file, json_file)

	if jsonFile != "" {
		json_file = "/vagrant/" + jsonFile
	}

	run_chef_solo := fmt.Sprintf("sudo chef-solo --config=%s --json-attributes=%s --override-runlist=%s --format=%s --log_level=%s",
		config_file, json_file, runlist, format, logLevel)

	cmd := strings.Join([]string{setup_dir, setup_config, setup_json, run_chef_solo}, " && ")
	return client.RunCommand(cmd)
}

func main() {
	log.SetFlags(0)

	var (
		host     = flag.String("H", "", "Set hostname for direct SSH access")
		machine  = flag.String("M", "", "Set name of Vagrant virtual machine")
		format   = flag.String("F", "null", "Set output format")
		logLevel = flag.String("l", "info", "Set log level")
		jsonFile = flag.String("j", "", "Load attributes from a JSON file")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: chef-runner [flags] [recipe ...]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	recipes := flag.Args()

	if *host != "" && *machine != "" {
		log.Fatal("error: -H and -M cannot be used together")
	}
	var client SSHClient
	if *host != "" {
		client = openssh.NewSSHClient(*host)
	} else {
		client = vagrant.NewSSHClient(*machine)
	}

	cb, err := cookbook.New(".")
	if err != nil {
		log.Fatal(err)
	}

	runlist := buildRunList(cb, recipes)
	fmt.Println("Run List is", runlist)

	if err := installCookbooks(cb, CookbookPath); err != nil {
		log.Fatal(err)
	}

	if err := provision(client, *format, *logLevel, *jsonFile, runlist); err != nil {
		log.Fatal(err)
	}
}
