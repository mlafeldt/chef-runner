package main

import (
	"flag"
	"fmt"
	"os"
)

var usage = `Usage: chef-runner [options] [--] [<recipe>...]

    -h              Show help text
    --version       Show program version

    -H <name>       Set hostname for direct SSH access
    -M <name>       Set name/UUID of Vagrant virtual machine

Options that will be passed to Chef Solo:

    -F <format>     Set output format (null, doc, minimal, min)
                    default: doc
    -l <level>      Set log level (debug, info, warn, error, fatal)
                    default: info
    -j <file>       Load attributes from a JSON file
`

// Flags stores the flags passed on the command line.
type Flags struct {
	host        string
	machine     string
	format      string
	logLevel    string
	jsonFile    string
	showVersion bool
}

// ParseFlags parses the command line and returns flags and recipes.
func ParseFlags(args []string) (*Flags, []string) {
	f := flag.NewFlagSet("chef-runner", flag.ExitOnError)
	f.Usage = func() { fmt.Fprintf(os.Stderr, usage) }

	var flags Flags
	f.StringVar(&flags.host, "H", "", "")
	f.StringVar(&flags.machine, "M", "", "")
	f.StringVar(&flags.format, "F", "", "")
	f.StringVar(&flags.logLevel, "l", "", "")
	f.StringVar(&flags.jsonFile, "j", "", "")
	f.BoolVar(&flags.showVersion, "version", false, "")

	f.Parse(args)
	return &flags, f.Args()
}
