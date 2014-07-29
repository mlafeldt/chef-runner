package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/cookbook"
	"github.com/mlafeldt/chef-runner/driver"
	"github.com/mlafeldt/chef-runner/driver/ssh"
	"github.com/mlafeldt/chef-runner/driver/vagrant"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner/util"
)

func logLevel() log.Level {
	l := log.LevelInfo
	e := os.Getenv("CHEF_RUNNER_LOG")
	if e == "" {
		return l
	}
	m := map[string]log.Level{
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
	text := `Usage: chef-runner [options] [--] [<recipe>...]

    -h              Show help text
    -H <name>       Set hostname for direct SSH access
    -M <name>       Set name of Vagrant virtual machine

Options that will be passed to Chef Solo:

    -F <format>     Set output format (null, doc, minimal, min)
                    default: null
    -l <level>      Set log level (debug, info, warn, error, fatal)
                    default: info
    -j <file>       Load attributes from a JSON file
`
	fmt.Fprintf(os.Stderr, text)
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
	log.SetLevel(logLevel())

	// usage() prints out flag documentation. No need to duplicate it here.
	var (
		host     = flag.String("H", "", "")
		machine  = flag.String("M", "", "")
		format   = flag.String("F", "", "")
		logLevel = flag.String("l", "", "")
		jsonFile = flag.String("j", "", "")
	)
	flag.Usage = usage
	flag.Parse()

	if *host != "" && *machine != "" {
		abort("-H and -M cannot be used together")
	}

	log.Info("Starting chef-runner")

	cb, err := cookbook.NewCookbook(".")
	if err != nil {
		abort(err)
	}
	if cb.Name == "" {
		abort("unknown cookbook name")
	}
	log.Debugf("Cookbook = %s\n", cb)

	recipes := flag.Args()
	runList := buildRunList(cb.Name, recipes)

	var attributes string
	if *jsonFile != "" {
		data, err := ioutil.ReadFile(*jsonFile)
		if err != nil {
			abort(err)
		}
		attributes = string(data)
	}

	var prov provisioner.Provisioner
	prov = chefsolo.Provisoner{
		RunList:    runList,
		Attributes: attributes,
		Format:     *format,
		LogLevel:   *logLevel,
	}

	log.Debugf("Provisoner = %+v\n", prov)

	if err := prov.CreateSandbox(); err != nil {
		abort(err)
	}

	// TODO: Copy files from p.SandboxPath to p.RootPath in order to get
	// rid of the Vagrant dependency

	var drv driver.Driver
	if *host != "" {
		drv, err = ssh.NewDriver(*host)
	} else {
		drv, err = vagrant.NewDriver(*machine)
	}
	if err != nil {
		abort(err)
	}

	log.Infof("Running Chef using %s\n", drv)

	cmd := strings.Join(prov.Command(), " ")
	log.Debug(cmd)

	if err := drv.RunCommand(cmd); err != nil {
		abort(err)
	}

	log.Info("chef-runner finished.")
}
