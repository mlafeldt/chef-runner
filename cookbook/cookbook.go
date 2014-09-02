// Package cookbook reads data from Chef cookbooks stored on disk. It can
// currently retrieve the cookbook's name and version as well as a list of all
// cookbook files.
package cookbook

import (
	"path"

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
func (cb Cookbook) Files() []string {
	fileList := [...]string{
		"README.md",
		"metadata.json",
		"metadata.rb",
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
	for _, f := range fileList {
		name := path.Join(cb.Path, f)
		if util.FileExist(name) {
			files = append(files, name)
		}
	}

	return files
}
