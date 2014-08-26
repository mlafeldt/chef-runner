package main

import (
	"flag"
	"fmt"
	"os"
)

var usage = `Usage: chef-runner [options] [--] [<recipe>...]

    -h, --help                   Show help text
    --version                    Show program version

    -H, --host <name>            Set hostname for direct SSH access
    -M, --machine <name>         Set name/UUID of Vagrant virtual machine

Options that will be passed to Chef Solo:

    -F, --format <format>        Set output format (null, doc, minimal, min)
                                 default: doc
    -l, --log_level <level>      Set log level (debug, info, warn, error, fatal)
                                 default: info
    -j, --json-attributes <file> Load attributes from a JSON file
`

// Flags stores the flags passed on the command line.
type Flags struct {
	Host        string
	Machine     string
	Format      string
	LogLevel    string
	JSONFile    string
	ShowVersion bool
}

// ParseFlags parses the command line and returns flags and recipes.
func ParseFlags(args []string) (*Flags, []string) {
	f := flag.NewFlagSet("chef-runner", flag.ExitOnError)
	f.Usage = func() { fmt.Fprintf(os.Stderr, usage) }

	var flags Flags

	f.BoolVar(&flags.ShowVersion, "version", false, "")

	f.StringVar(&flags.Host, "H", "", "")
	f.StringVar(&flags.Host, "host", "", "")

	f.StringVar(&flags.Machine, "M", "", "")
	f.StringVar(&flags.Machine, "machine", "", "")

	f.StringVar(&flags.Format, "F", "", "")
	f.StringVar(&flags.Format, "format", "", "")

	f.StringVar(&flags.LogLevel, "l", "", "")
	f.StringVar(&flags.LogLevel, "log_level", "", "")

	f.StringVar(&flags.JSONFile, "j", "", "")
	f.StringVar(&flags.JSONFile, "json-attributes", "", "")

	f.Parse(args)
	return &flags, f.Args()
}
