package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var usage = `Usage: chef-runner [options] [--] [<recipe>...]

  -H, --host <name>            Name of host reachable over SSH
  -M, --machine <name>         Name or UUID of Vagrant virtual machine
  -K, --kitchen <name>         Name of Test Kitchen instance

  --ssh-option <key=value>     Specify custom SSH option, can be used multiple times

  -i, --install-chef <version> Install Chef (x.y.z, latest, true, false)
                               default: false

  -F, --format <format>        Chef output format (null, doc, minimal, min)
                               default: doc
  -l, --log_level <level>      Chef log level (debug, info, warn, error, fatal)
                               default: info
  -j, --json-attributes <file> Load attributes from a JSON file

  -h, --help                   Show help text
  --version                    Show program version
`

// Flags stores the options and arguments passed on the command line.
type Flags struct {
	Host    string
	Machine string
	Kitchen string

	SSHOptions map[string]string

	Format   string
	LogLevel string
	JSONFile string

	ShowVersion bool
	ChefVersion string

	Recipes []string
}

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// ParseFlags parses the command line and returns the result.
func ParseFlags(args []string) (*Flags, error) {
	f := flag.NewFlagSet("chef-runner", flag.ExitOnError)
	f.Usage = func() { fmt.Fprintf(os.Stderr, usage) }

	var flags Flags

	f.BoolVar(&flags.ShowVersion, "version", false, "")

	f.StringVar(&flags.Host, "H", "", "")
	f.StringVar(&flags.Host, "host", "", "")

	f.StringVar(&flags.Machine, "M", "", "")
	f.StringVar(&flags.Machine, "machine", "", "")

	f.StringVar(&flags.Kitchen, "K", "", "")
	f.StringVar(&flags.Kitchen, "kitchen", "", "")

	var sshOptions stringSlice
	f.Var(&sshOptions, "ssh-option", "")

	f.StringVar(&flags.Format, "F", "", "")
	f.StringVar(&flags.Format, "format", "", "")

	f.StringVar(&flags.LogLevel, "l", "", "")
	f.StringVar(&flags.LogLevel, "log_level", "", "")

	f.StringVar(&flags.JSONFile, "j", "", "")
	f.StringVar(&flags.JSONFile, "json-attributes", "", "")

	f.StringVar(&flags.ChefVersion, "i", "", "")
	f.StringVar(&flags.ChefVersion, "install-chef", "", "")

	if err := f.Parse(args); err != nil {
		return nil, err
	}

	if flags.Host != "" && flags.Machine != "" {
		return nil, errors.New("-H and -M cannot be used together")
	}

	if len(sshOptions) > 0 {
		flags.SSHOptions = make(map[string]string)
		for _, o := range sshOptions {
			fields := strings.Split(o, "=")
			if len(fields) != 2 {
				return nil, errors.New("invalid SSH option: " + o)
			}
			flags.SSHOptions[fields[0]] = fields[1]
		}
	}

	if len(f.Args()) > 0 {
		flags.Recipes = f.Args()
	}

	return &flags, nil
}
