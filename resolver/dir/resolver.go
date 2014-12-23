// Package dir implements a cookbook dependency resolver that merely copies
// cookbook directories to the right place.
package dir

import (
	"errors"
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/chef/cookbook"
	"github.com/mlafeldt/chef-runner/rsync"
)

// Resolver is a cookbook dependency resolver that copies cookbook directories
// to the right place.
type Resolver struct{}

func installCookbook(dst, src string) error {
	cb, err := cookbook.NewCookbook(src)
	if err != nil {
		return err
	}

	if cb.Name == "" {
		return errors.New("cookbook name required")
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	return rsync.MirrorClient.Copy(path.Join(dst, cb.Name), cb.Files()...)
}

// Resolve copies the cookbook in the current directory to dst.
func (r Resolver) Resolve(dst string) error {
	return installCookbook(dst, ".")
}

// Name returns the resolver's name.
func (r Resolver) Name() string {
	return "Directory"
}
