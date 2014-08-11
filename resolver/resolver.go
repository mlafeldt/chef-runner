// Package resolver provides a generic cookbook dependency resolver.
package resolver

import (
	"errors"

	"github.com/mlafeldt/chef-runner/cookbook"
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
	String() string
}

// Helper to determine resolver from files in current directory.
func findResolver(dst string) (Resolver, error) {
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

// AutoResolve automatically resolves cookbook dependencies based on the files
// present in the current directory.
func AutoResolve(dst string) error {
	r, err := findResolver(dst)
	if err != nil {
		return err
	}

	log.Infof("Installing cookbook dependencies with %s\n", r)
	return r.Resolve(dst)
}
