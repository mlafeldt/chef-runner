// Package cli handles the command line interface of chef-runner. This includes
// parsing of options and arguments as well as printing help text.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
)

var usage = `Usage: chef-runner [options] [--] [<recipe>...]

  -H, --host <name>            Name of host reachable over SSH
  -M, --machine <name>         Name or UUID of Vagrant virtual machine
  -K, --kitchen <name>         Name of Test Kitchen instance

  --ssh-option <option>        Add OpenSSH option as specified in ssh_config(5)
  --rsync-option <option>      Add Rsync option as listed in rsync(1)

  -i, --install-chef <version> Install Chef (x.y.z, latest, true, false)
                               default: false

  -F, --format <format>        Chef output format (null, doc, minimal, min)
                               default: doc
  -l, --log_level <level>      Chef log level (debug, info, warn, error, fatal)
                               default: info
  -j, --json-attributes <file> Load attributes from a JSON file

  --color=false                Disable colorized output (enabled by default)
  -h, --help                   Show help text
  --version                    Show program version
`

// This slice is used to implement options that can be passed multiple times.
type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Flags stores the options and arguments passed on the command line.
type Flags struct {
	Host    string
	Machine string
	Kitchen string

	SSHOptions   stringSlice
	RsyncOptions stringSlice

	ChefVersion string
	Format      string
	LogLevel    string
	JSONFile    string
	Recipes     []string

	Color       bool
	ShowVersion bool
}

// ParseFlags parses the command line and returns the result.
func ParseFlags(args []string) (*Flags, error) {
	f := flag.NewFlagSet("chef-runner", flag.ExitOnError)
	f.Usage = func() { fmt.Fprintf(os.Stderr, usage) }

	var flags Flags

	f.StringVar(&flags.Host, "H", "", "")
	f.StringVar(&flags.Host, "host", "", "")

	f.StringVar(&flags.Machine, "M", "", "")
	f.StringVar(&flags.Machine, "machine", "", "")

	f.StringVar(&flags.Kitchen, "K", "", "")
	f.StringVar(&flags.Kitchen, "kitchen", "", "")

	f.Var(&flags.SSHOptions, "ssh-option", "")

	f.Var(&flags.RsyncOptions, "rsync-option", "")

	f.StringVar(&flags.ChefVersion, "i", "", "")
	f.StringVar(&flags.ChefVersion, "install-chef", "", "")

	f.StringVar(&flags.Format, "F", "", "")
	f.StringVar(&flags.Format, "format", "", "")

	f.StringVar(&flags.LogLevel, "l", "", "")
	f.StringVar(&flags.LogLevel, "log_level", "", "")

	f.StringVar(&flags.JSONFile, "j", "", "")
	f.StringVar(&flags.JSONFile, "json-attributes", "", "")

	f.BoolVar(&flags.Color, "color", true, "")

	f.BoolVar(&flags.ShowVersion, "version", false, "")

	if err := f.Parse(args); err != nil {
		return nil, err
	}

	n := 0
	for _, i := range []string{flags.Host, flags.Machine, flags.Kitchen} {
		if i != "" {
			n++
		}
	}
	if n > 1 {
		return nil, errors.New("-H, -M, and -K cannot be used together")
	}

	if len(f.Args()) > 0 {
		flags.Recipes = f.Args()
	}

	return &flags, nil
}

// LogLevel returns the log level to use based on the CHEF_RUNNER_LOG
// environment variable.
func LogLevel() log.Level {
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
