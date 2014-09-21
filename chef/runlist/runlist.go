// Package runlist builds Chef run lists. chef-runner allows to compose run
// lists using a flexible recipe syntax. If required, this package translates
// that syntax to Chef's syntax.
package runlist

import (
	"errors"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/util"
)

func expand(recipe, cookbook string) (string, error) {
	if strings.HasPrefix(recipe, "::") {
		if cookbook == "" {
			log.Errorf("cannot add local recipe \"%s\" to run list\n",
				strings.TrimPrefix(recipe, "::"))
			return "", errors.New("cookbook name required")
		}
		return cookbook + recipe, nil
	}
	if path.Dir(recipe) == "recipes" && path.Ext(recipe) == ".rb" {
		if cookbook == "" {
			log.Errorf("cannot add local recipe \"%s\" to run list\n", recipe)
			return "", errors.New("cookbook name required")
		}
		return cookbook + "::" + util.BaseName(recipe, ".rb"), nil
	}
	return recipe, nil
}

// Build creates a Chef run list from a list of recipes and an optional
// cookbook name. The cookbook name is only required to expand local recipes.
func Build(recipes []string, cookbook string) ([]string, error) {
	runList := []string{}
	for _, r := range recipes {
		for _, r := range strings.Split(r, ",") {
			recipe, err := expand(r, cookbook)
			if err != nil {
				return nil, err
			}
			runList = append(runList, recipe)
		}
	}
	return runList, nil
}
