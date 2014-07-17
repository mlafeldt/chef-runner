package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner.go/cookbook"
	"github.com/mlafeldt/chef-runner.go/log"
	"github.com/mlafeldt/chef-runner.go/openssh"
	"github.com/mlafeldt/chef-runner.go/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner.go/util"
	"github.com/mlafeldt/chef-runner.go/vagrant"
)

type SSHClient interface {
	RunCommand(command string) error
	String() string
}

func logLevel() int {
	l := log.LevelInfo
	e := os.Getenv("CHEF_RUNNER_LOG")
	if e == "" {
		return l
	}
	m := map[string]int{
		"debug": log.LevelDebug,
		"info":  log.LevelInfo,
		"warn":  log.LevelWarn,
		"error": log.LevelError,
	}
	if v, ok := m[strings.ToLower(e)]; ok {
		l = v
	}
	return l
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: chef-runner [flags] [recipe ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func abort(v ...interface{}) {
	log.Error(v...)
	os.Exit(1)
}

func buildRunList(cookbookName string, recipes []string) []string {
	if len(recipes) == 0 {
		return []string{cookbookName + "::default"}
	}

	var runList []string
	for _, r := range recipes {
		var recipeName string
		if strings.Contains(r, "::") {
			recipeName = r
		} else if path.Dir(r) == "recipes" && path.Ext(r) == ".rb" {
			recipeName = cookbookName + "::" + util.BaseName(r, ".rb")
		} else {
			recipeName = cookbookName + "::" + r
		}
		runList = append(runList, recipeName)
	}
	return runList
}

func main() {
	log.Level = logLevel()

	var (
		host     = flag.String("H", "", "Set hostname for direct SSH access")
		machine  = flag.String("M", "", "Set name of Vagrant virtual machine")
		format   = flag.String("F", "", "Set output format")
		logLevel = flag.String("l", "", "Set log level")
		jsonFile = flag.String("j", "", "Load attributes from a JSON file")
	)
	flag.Usage = usage
	flag.Parse()

	if *host != "" && *machine != "" {
		abort("-H and -M cannot be used together")
	}

	cb, err := cookbook.NewCookbook(".")
	if err != nil {
		abort(err)
	}
	if cb.Name == "" {
		abort("unknown cookbook name")
	}
	log.Debug("Cookbook =", cb)

	recipes := flag.Args()
	runList := buildRunList(cb.Name, recipes)
	log.Debug("Run list =", runList)

	var attributes string
	if *jsonFile != "" {
		data, err := ioutil.ReadFile(*jsonFile)
		if err != nil {
			abort(err)
		}
		attributes = string(data)
	}
	log.Debug("Attributes =", attributes)

	provisioner := chefsolo.Provisoner{
		RunList:    runList,
		Attributes: attributes,
		Format:     *format,
		LogLevel:   *logLevel,
	}
	if err := provisioner.CreateSandbox(); err != nil {
		abort(err)
	}

	// TODO: Copy files from p.SandboxPath to p.RootPath in order to get
	// rid of the Vagrant dependency

	var client SSHClient
	if *host != "" {
		client = openssh.NewClient(*host)
	} else {
		client = vagrant.NewClient(*machine)
	}

	log.Info("Running Chef using " + client.String())
	cmd := strings.Join(provisioner.Command(), " ")
	log.Debug(cmd)

	if err := client.RunCommand(cmd); err != nil {
		abort(err)
	}
}
