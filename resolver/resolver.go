// Package resolver provides a generic cookbook dependency resolver.
package resolver

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/chef/cookbook"
	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/resolver/berkshelf"
	"github.com/mlafeldt/chef-runner/resolver/dir"
	"github.com/mlafeldt/chef-runner/resolver/librarian"
	"github.com/mlafeldt/chef-runner/util"
)

// A Resolver resolves cookbook dependencies and installs them to directory dst.
// This is the interface that all resolvers need to implement.
type Resolver interface {
	Resolve(dst string) error
	Name() string
	InputFiles() []string
}

func findResolverByName(name string) (Resolver, error) {
	var resolvers = [...]Resolver{
		berkshelf.Resolver{},
		librarian.Resolver{},
		dir.Resolver{},
	}
	for _, r := range resolvers {
		if strings.HasPrefix(strings.ToLower(r.Name()), strings.ToLower(name)) {
			return r, nil
		}
	}
	return nil, errors.New("unknown resolver name: " + name)
}

func findResolver(name, dst string) (Resolver, error) {
	if name != "" {
		return findResolverByName(name)
	}

	tsFile := path.Join(dst, "action_resolve")
	ts, _ := util.ReadTimestampFile(tsFile)

	for _, name := range berkshelf.InputFiles {
		mt, err := util.FileModTime(name)
		if err == nil && mt >= ts {
			log.Debugf("%s was updated, using Berkshelf\n", name)
			return berkshelf.Resolver{}, nil
		}
	}

	for _, name := range librarian.InputFiles {
		mt, err := util.FileModTime(name)
		if err == nil && mt >= ts {
			log.Debugf("%s was updated, using Librarian\n", name)
			return librarian.Resolver{}, nil
		}
	}

	cb, _ := cookbook.NewCookbook(".")
	if cb.Name != "" {
		return dir.Resolver{}, nil
	}

	log.Error("Berksfile, Cheffile, or metadata.rb must exist in current directory")
	return nil, errors.New("cookbooks could not be found")
}

func stripCookbooks(dst string) error {
	cookbookDirs, err := ioutil.ReadDir(dst)
	if err != nil {
		return err
	}

	for _, dir := range cookbookDirs {
		if !dir.IsDir() {
			continue
		}
		cb := cookbook.Cookbook{Path: path.Join(dst, dir.Name())}
		if err := cb.Strip(); err != nil {
			return err
		}
	}

	return nil
}

// Resolve resolves cookbook dependencies using the named resolver. If no
// resolver is specified, it will be guessed based on the files present in the
// current directory. After resolving dependencies, all non-cookbook files will
// be deleted as well.
func Resolve(name, dst string) error {
	log.Debug("Preparing cookbooks")

	r, err := findResolver(name, dst)
	if err != nil {
		return err
	}

	log.Infof("Installing cookbook dependencies with %s resolver\n", r.Name())
	if err := r.Resolve(dst); err != nil {
		return err
	}

	tsFile := path.Join(dst, "action_resolve")
	if err := util.WriteTimestampFile(tsFile); err != nil {
		return err
	}

	log.Info("Stripping non-cookbook files")
	return stripCookbooks(dst)

}

// AutoResolve automatically resolves cookbook dependencies based on the files
// present in the current directory. After resolving dependencies, all
// non-cookbook files will be deleted as well.
func AutoResolve(dst string) error {
	return Resolve("", dst)
}
