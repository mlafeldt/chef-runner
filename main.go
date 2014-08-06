package main

import (
	"errors"
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
    --version       Show program version

    -H <name>       Set hostname for direct SSH access
    -M <name>       Set name/UUID of Vagrant virtual machine

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

func buildRunList(cookbookName string, recipes []string) ([]string, error) {
	if len(recipes) == 0 {
		recipes = []string{"default"}
	}

	var runList []string
	for _, r := range recipes {
		var recipeName string
		if strings.Contains(r, "::") {
			recipeName = r
		} else {
			if cookbookName == "" {
				log.Errorf("cannot add recipe `%s` to run list\n", r)
				return nil, errors.New("cookbook name required")
			}
			if path.Dir(r) == "recipes" && path.Ext(r) == ".rb" {
				recipeName = cookbookName + "::" + util.BaseName(r, ".rb")
			} else {
				recipeName = cookbookName + "::" + r
			}
		}
		runList = append(runList, recipeName)
	}
	return runList, nil
}

func main() {
	log.SetLevel(logLevel())

	// usage() prints out flag documentation. No need to duplicate it here.
	var (
		host        = flag.String("H", "", "")
		machine     = flag.String("M", "", "")
		format      = flag.String("F", "", "")
		logLevel    = flag.String("l", "", "")
		jsonFile    = flag.String("j", "", "")
		showVersion = flag.Bool("version", false, "")
	)
	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("chef-runner %s %s\n", VersionString(), TargetString())
		os.Exit(0)
	}

	if *host != "" && *machine != "" {
		abort("-H and -M cannot be used together")
	}

	log.Infof("Starting chef-runner (%s %s)\n", VersionString(), TargetString())

	cb, err := cookbook.NewCookbook(".")
	if err != nil {
		abort(err)
	}
	log.Debugf("Cookbook = %s\n", cb)

	runList, err := buildRunList(cb.Name, flag.Args())
	if err != nil {
		abort(err)
	}

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

	var drv driver.Driver
	if *host != "" {
		drv, err = ssh.NewDriver(*host)
	} else {
		drv, err = vagrant.NewDriver(*machine)
	}
	if err != nil {
		abort(err)
	}

	log.Info("Uploading local files to machine. This may take a while...")
	log.Debugf("Uploading files from %s to %s on machine\n",
		provisioner.SandboxPath, provisioner.RootPath)
	if err := drv.Upload(provisioner.RootPath, provisioner.SandboxPath+"/"); err != nil {
		abort(err)
	}

	log.Infof("Running Chef using %s\n", drv)
	cmd := strings.Join(prov.Command(), " ")
	if err := drv.RunCommand(cmd); err != nil {
		abort(err)
	}

	log.Info("chef-runner finished.")
}
