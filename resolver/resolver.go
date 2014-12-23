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

	cb, _ := cookbook.NewCookbook(".")

	// If the current folder is a cookbook and its dependencies have
	// already been resolved, only update this cookbook with rsync.
	// TODO: improve this check by comparing timestamps etc.
	if cb.Name != "" && util.FileExist(dst) {
		return dir.Resolver{}, nil
	}

	if util.FileExist("Berksfile") {
		return berkshelf.Resolver{}, nil
	}

	if util.FileExist("Cheffile") {
		return librarian.Resolver{}, nil
	}

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

	log.Info("Stripping non-cookbook files")
	return stripCookbooks(dst)

}

// AutoResolve automatically resolves cookbook dependencies based on the files
// present in the current directory. After resolving dependencies, all
// non-cookbook files will be deleted as well.
func AutoResolve(dst string) error {
	return Resolve("", dst)
}
