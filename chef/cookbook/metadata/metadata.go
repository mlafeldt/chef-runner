// Package metadata parses Chef cookbook metadata. It can currently retrieve
// the cookbook's name and version.
package metadata

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// Filename is the name of the cookbook file that stores metadata.
const Filename = "metadata.rb"

// Metadata stores metadata about a cookbook.
type Metadata struct {
	Name    string
	Version string
}

// Parse parses cookbook metadata from an io.Reader. It returns Metadata.
func Parse(r io.Reader) (*Metadata, error) {
	metadata := Metadata{}
	scanner := bufio.NewScanner(r)
	re := regexp.MustCompile(`\A(\S+)\s+['"](.*?)['"]\z`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		match := re.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		switch match[1] {
		case "name":
			metadata.Name = match[2]
		case "version":
			metadata.Version = match[2]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// ParseFile parses a cookbook metadata file. It returns Metadata.
func ParseFile(name string) (*Metadata, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}
