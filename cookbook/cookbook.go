package cookbook

import (
	"path"
	"path/filepath"

	"github.com/mlafeldt/chef-runner/cookbook/metadata"
	"github.com/mlafeldt/chef-runner/util"
)

type Cookbook struct {
	Path    string
	Name    string
	Version string
}

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

func (cb Cookbook) String() string {
	// TODO: check if fields are actually set
	return cb.Name + " " + cb.Version
}

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
