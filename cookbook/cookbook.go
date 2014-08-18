// Package cookbook reads data from Chef cookbooks stored on disk. It can
// currently retrieve the cookbook's name and version as well as a list of all
// cookbook files.
package cookbook

import (
	"path"
	"path/filepath"

	"github.com/mlafeldt/chef-runner/cookbook/metadata"
	"github.com/mlafeldt/chef-runner/util"
)

// A Cookbook is a Chef cookbook stored on disk.
type Cookbook struct {
	Path    string
	Name    string
	Version string
}

// NewCookbook returns a Cookbook that is located at cookbookPath.
func NewCookbook(cookbookPath string) (*Cookbook, error) {
	cb := Cookbook{Path: cookbookPath}

	metadataPath := path.Join(cookbookPath, metadata.Filename)
	if util.FileExist(metadataPath) {
		metadata, err := metadata.ParseFile(metadataPath)
		if err != nil {
			return nil, err
		}
		cb.Name = metadata.Name
		cb.Version = metadata.Version
	}

	return &cb, nil
}

// String returns the cookbook's name and version.
func (cb Cookbook) String() string {
	// TODO: check if fields are actually set
	return cb.Name + " " + cb.Version
}

// Files returns the names of all cookbook files. Other files are ignored.
func (cb Cookbook) Files() ([]string, error) {
	filesGlob := []string{
		"README.*",
		"metadata.*",
		"attributes",
		"definitions",
		"files",
		"libraries",
		"providers",
		"recipes",
		"resources",
		"templates",
	}
	var files []string
	for _, glob := range filesGlob {
		matches, err := filepath.Glob(path.Join(cb.Path, glob))
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}
	return files, nil
}
