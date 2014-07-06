package cookbook

import (
	"path"
	"path/filepath"
	"strings"
)

func Files(cookbookPath string) ([]string, error) {
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
		matches, err := filepath.Glob(path.Join(cookbookPath, glob))
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}
	return files, nil
}

func NameFromPath(cookbookPath string) string {
	base := path.Base(cookbookPath)
	if strings.HasPrefix(base, "chef-") {
		return strings.TrimPrefix(base, "chef-")
	}
	if strings.HasSuffix(base, "-cookbook") {
		return strings.TrimSuffix(base, "-cookbook")
	}
	return base
}
