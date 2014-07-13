package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner.go/cookbook"
	"github.com/mlafeldt/chef-runner.go/openssh"
	"github.com/mlafeldt/chef-runner.go/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner.go/util"
	"github.com/mlafeldt/chef-runner.go/vagrant"
)

const (
	CookbookPath = "vendor/cookbooks"
)

type SSHClient interface {
	RunCommand(command string) error
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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: chef-runner [flags] [recipe ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)

	var (
		host     = flag.String("H", "", "Set hostname for direct SSH access")
		machine  = flag.String("M", "", "Set name of Vagrant virtual machine")
		format   = flag.String("F", chefsolo.DefaultFormat, "Set output format")
		logLevel = flag.String("l", chefsolo.DefaultLogLevel, "Set log level")
		jsonFile = flag.String("j", "", "Load attributes from a JSON file")
	)
	flag.Usage = usage
	flag.Parse()

	if *host != "" && *machine != "" {
		log.Fatal("error: -H and -M cannot be used together")
	}
	var client SSHClient
	if *host != "" {
		client = openssh.NewClient(*host)
	} else {
		client = vagrant.NewClient(*machine)
	}

	cb, err := cookbook.New(".")
	if err != nil {
		log.Fatal(err)
	}

	recipes := flag.Args()
	runList := buildRunList(cb.Name, recipes)
	log.Println("Run List is", runList)

	var attributes string
	if *jsonFile != "" {
		data, err := ioutil.ReadFile(*jsonFile)
		if err != nil {
			log.Fatal(err)
		}
		attributes = string(data)
	}

	p := chefsolo.Provisoner{
		RunList:    runList,
		Attributes: attributes,
		Format:     *format,
		LogLevel:   *logLevel,
	}
	if err := p.CreateSandbox(); err != nil {
		log.Fatal(err)
	}
	cmd := strings.Join(p.Command(), " ")
	log.Println(cmd)
	if err := client.RunCommand(cmd); err != nil {
		log.Fatal(err)
	}
}
