package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mlafeldt/chef-runner/chef/cookbook"
	"github.com/mlafeldt/chef-runner/chef/omnibus"
	"github.com/mlafeldt/chef-runner/chef/runlist"
	"github.com/mlafeldt/chef-runner/cli"
	"github.com/mlafeldt/chef-runner/driver"
	"github.com/mlafeldt/chef-runner/driver/kitchen"
	"github.com/mlafeldt/chef-runner/driver/ssh"
	"github.com/mlafeldt/chef-runner/driver/vagrant"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/provisioner"
	"github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner/resolver"
)

const (
	// SandboxPath is the path to the local sandbox directory where
	// chef-runner stores files that will be uploaded to a machine.
	SandboxPath = ".chef-runner/sandbox"

	// RootPath is the path on the machine where files from SandboxPath
	// will be uploaded to.
	RootPath = "/tmp/chef-runner"
)

func abort(v ...interface{}) {
	log.Error(v...)
	os.Exit(1)
}

func findDriver(flags *cli.Flags) (driver.Driver, error) {
	if flags.Host != "" {
		return ssh.NewDriver(flags.Host, flags.SSHOptions, flags.RsyncOptions)
	}
	if flags.Kitchen != "" {
		return kitchen.NewDriver(flags.Kitchen, flags.SSHOptions, flags.RsyncOptions)
	}
	return vagrant.NewDriver(flags.Machine, flags.SSHOptions, flags.RsyncOptions)
}

func uploadFiles(drv driver.Driver) error {
	log.Info("Uploading local files to machine. This may take a while...")
	log.Debugf("Uploading files from %s to %s on machine\n",
		SandboxPath, RootPath)
	return drv.Upload(RootPath, SandboxPath+"/")
}

func installChef(drv driver.Driver, installer omnibus.Installer) error {
	installCmd := installer.Command()
	if len(installCmd) == 0 {
		log.Info("Skipping installation of Chef")
		return nil
	}
	log.Info("Installing Chef")
	return drv.RunCommand(installCmd)
}

func runChef(drv driver.Driver, p provisioner.Provisioner) error {
	log.Infof("Running Chef using %s\n", drv)
	return drv.RunCommand(p.Command())
}

func main() {
	startTime := time.Now()

	log.SetLevel(cli.LogLevel())

	flags, err := cli.ParseFlags(os.Args[1:])
	if err != nil {
		abort(err)
	}

	log.UseColor = flags.Color

	if flags.ShowVersion {
		fmt.Printf("chef-runner %s %s %s\n", VersionString(),
			TargetString(), GoVersionString())
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
		RunList:    runList,
		Attributes: attributes,
		Format:     flags.Format,
		LogLevel:   flags.LogLevel,
		UseSudo:    true,

		SandboxPath: SandboxPath,
		RootPath:    RootPath,
	}
	log.Debugf("Provisioner = %+v\n", p)

	installer := omnibus.Installer{
		ChefVersion: flags.ChefVersion,
		SandboxPath: SandboxPath,
		RootPath:    RootPath,
	}
	log.Debugf("Installer = %+v\n", installer)

	log.Info("Preparing local files")

	log.Debug("Creating local sandbox in", SandboxPath)
	if err := os.MkdirAll(SandboxPath, 0755); err != nil {
		abort(err)
	}

	if err := p.PrepareFiles(); err != nil {
		abort(err)
	}

	if err := resolver.AutoResolve(path.Join(SandboxPath, "cookbooks")); err != nil {
		abort(err)
	}

	if err := installer.PrepareFiles(); err != nil {
		abort(err)
	}

	drv, err := findDriver(flags)
	if err != nil {
		abort(err)
	}

	if err := uploadFiles(drv); err != nil {
		abort(err)
	}

	if err := installChef(drv, installer); err != nil {
		abort(err)
	}

	if err := runChef(drv, p); err != nil {
		abort(err)
	}

	log.Info("chef-runner finished in", time.Now().Sub(startTime))
}
