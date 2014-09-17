// Package cookbook reads and manipulates data from Chef cookbooks stored on
// disk.
package cookbook

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/chef/cookbook/metadata"
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

// Strip removes all non-cookbook files from the cookbook.
func (cb Cookbook) Strip() error {
	cbFiles := make(map[string]bool)
	for _, f := range cb.Files() {
		cbFiles[f] = true
	}

	files, err := ioutil.ReadDir(cb.Path)
	if err != nil {
		return err
	}

	for _, f := range files {
		name := path.Join(cb.Path, f.Name())
		if _, keep := cbFiles[name]; !keep {
			if err := os.RemoveAll(name); err != nil {
				return err
			}
		}
	}

	return nil
}
