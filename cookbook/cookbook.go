package cookbook

import (
	"path"
	"path/filepath"

	"github.com/mlafeldt/chef-runner.go/cookbook/metadata"
)

type Cookbook struct {
	Path    string
	Name    string
	Version string
}

func NewCookbook(cookbookPath string) (*Cookbook, error) {
	metadata, err := metadata.ParseFile(path.Join(cookbookPath, "metadata.rb"))
	if err != nil {
		return nil, err
	}
	cb := Cookbook{
		Path:    cookbookPath,
		Name:    metadata.Name,
		Version: metadata.Version,
	}
	return &cb, nil
}

func (cb *Cookbook) Files() ([]string, error) {
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
