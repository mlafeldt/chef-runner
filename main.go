package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/mlafeldt/chef-runner/chef/cookbook"
	"github.com/mlafeldt/chef-runner/chef/runlist"
	"github.com/mlafeldt/chef-runner/driver"
	"github.com/mlafeldt/chef-runner/driver/kitchen"
	"github.com/mlafeldt/chef-runner/driver/ssh"
	"github.com/mlafeldt/chef-runner/driver/vagrant"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/provisioner/chefsolo"
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

func findDriver(flags *Flags) (driver.Driver, error) {
	if flags.Host != "" {
		return ssh.NewDriver(flags.Host, nil)
	}
	if flags.Kitchen != "" {
		return kitchen.NewDriver(flags.Kitchen, nil)
	}
	return vagrant.NewDriver(flags.Machine, nil)
}

func uploadFiles(drv driver.Driver) error {
	log.Info("Uploading local files to machine. This may take a while...")
	log.Debugf("Uploading files from %s to %s on machine\n",
		provisioner.SandboxPath, provisioner.RootPath)
	return drv.Upload(provisioner.RootPath, provisioner.SandboxPath+"/")
}

func installChef(drv driver.Driver, p provisioner.Provisioner) error {
	installCmd := p.InstallCommand()
	if len(installCmd) == 0 {
		log.Info("Skipping installation of Chef")
		return nil
	}
	log.Info("Installing Chef")
	return drv.RunCommand(installCmd)
}

func runChef(drv driver.Driver, p provisioner.Provisioner) error {
	log.Infof("Running Chef using %s\n", drv)
	return drv.RunCommand(p.ProvisionCommand())
}

func main() {
	startTime := time.Now()

	log.SetLevel(logLevel())

	flags, err := ParseFlags(os.Args[1:])
	if err != nil {
		abort(err)
	}

	if flags.ShowVersion {
		fmt.Printf("chef-runner %s %s\n", VersionString(), TargetString())
		os.Exit(0)
	}

	log.Infof("Starting chef-runner (%s %s)\n", VersionString(), TargetString())

	var attributes string
	if flags.JSONFile != "" {
		data, err := ioutil.ReadFile(flags.JSONFile)
		if err != nil {
			abort(err)
		}
		attributes = string(data)
	}

	// 1) Run default recipe if no recipes are passed
	// 2) Use run list from JSON file if present, overriding 1)
	// 3) Use run list from command line if present, overriding 1) and 2)
	recipes := flags.Recipes
	if len(recipes) == 0 {
		// TODO: parse actual JSON data
		if strings.Contains(attributes, "run_list") {
			log.Infof("Using run list from %s\n", flags.JSONFile)
		} else {
			recipes = []string{"::default"}
		}
	}

	var runList []string
	if len(recipes) > 0 {
		cb, err := cookbook.NewCookbook(".")
		if err != nil {
			abort(err)
		}
		log.Debugf("Cookbook = %s\n", cb)

		if runList, err = runlist.Build(recipes, cb.Name); err != nil {
			abort(err)
		}
		log.Infof("Run list is %s\n", runList)
	}

	var p provisioner.Provisioner
	p = chefsolo.Provisioner{
		RunList:     runList,
		Attributes:  attributes,
		Format:      flags.Format,
		LogLevel:    flags.LogLevel,
		UseSudo:     true,
		ChefVersion: flags.ChefVersion,
	}

	log.Debugf("Provisioner = %+v\n", p)

	if err := p.CreateSandbox(); err != nil {
		abort(err)
	}

	drv, err := findDriver(flags)
	if err != nil {
		abort(err)
	}

	if err := uploadFiles(drv); err != nil {
		abort(err)
	}

	if err := installChef(drv, p); err != nil {
		abort(err)
	}

	if err := runChef(drv, p); err != nil {
		abort(err)
	}

	log.Info("chef-runner finished in", time.Now().Sub(startTime))
}
