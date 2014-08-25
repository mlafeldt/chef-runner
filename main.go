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

	var flags struct {
		host        string
		machine     string
		format      string
		logLevel    string
		jsonFile    string
		showVersion bool
	}
	// Usage prints out flag documentation. No need to duplicate it here.
	flag.Usage = Usage
	flag.StringVar(&flags.host, "H", "", "")
	flag.StringVar(&flags.machine, "M", "", "")
	flag.StringVar(&flags.format, "F", "", "")
	flag.StringVar(&flags.logLevel, "l", "", "")
	flag.StringVar(&flags.jsonFile, "j", "", "")
	flag.BoolVar(&flags.showVersion, "version", false, "")
	flag.Parse()

	if flags.showVersion {
		fmt.Printf("chef-runner %s %s\n", VersionString(), TargetString())
		os.Exit(0)
	}

	if flags.host != "" && flags.machine != "" {
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
	if flags.jsonFile != "" {
		data, err := ioutil.ReadFile(flags.jsonFile)
		if err != nil {
			abort(err)
		}
		attributes = string(data)
	}

	var prov provisioner.Provisioner
	prov = chefsolo.Provisioner{
		RunList:    runList,
		Attributes: attributes,
		Format:     flags.format,
		LogLevel:   flags.logLevel,
		UseSudo:    true,
	}

	log.Debugf("Provisioner = %+v\n", prov)

	if err := prov.CreateSandbox(); err != nil {
		abort(err)
	}

	var drv driver.Driver
	if flags.host != "" {
		drv, err = ssh.NewDriver(flags.host)
	} else {
		drv, err = vagrant.NewDriver(flags.machine)
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
	if err := drv.RunCommand(prov.Command()); err != nil {
		abort(err)
	}

	log.Info("chef-runner finished.")
}
