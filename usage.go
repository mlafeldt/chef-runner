package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var usage = `Usage: chef-runner [options] [--] [<recipe>...]

    -h, --help                   Show help text
    --version                    Show program version

    -H, --host <name>            Set hostname for direct SSH access
    -M, --machine <name>         Set name/UUID of Vagrant virtual machine
    -K, --kitchen <name>         Set name of Test Kitchen instance

    -i, --install-chef <version> Install Chef (x.y.z, latest, true, false)
                                 default: false

Options that will be passed to Chef Solo:

    -F, --format <format>        Set output format (null, doc, minimal, min)
                                 default: doc
    -l, --log_level <level>      Set log level (debug, info, warn, error, fatal)
                                 default: info
    -j, --json-attributes <file> Load attributes from a JSON file
`

// Flags stores the options and arguments passed on the command line.
type Flags struct {
	Host    string
	Machine string
	Kitchen string

	Format      string
	LogLevel    string
	JSONFile    string
	ShowVersion bool
	ChefVersion string

	Recipes []string
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

	if len(f.Args()) > 0 {
		flags.Recipes = f.Args()
	}

	return &flags, nil
}
