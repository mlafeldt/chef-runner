package metadata

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

const (
	Filename = "metadata.rb"
)

type Metadata struct {
	Name    string
	Version string
}

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

func ParseFile(name string) (*Metadata, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}
