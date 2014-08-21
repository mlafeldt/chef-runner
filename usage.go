package main

import (
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

// Usage prints the usage text to stderr. It is called by the flag package when
// either an invalid flag or -h is passed on the command line.
func Usage() {
	fmt.Fprintf(os.Stderr, usage)
}
