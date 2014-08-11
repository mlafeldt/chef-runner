// Package dir implements a cookbook dependency resolver that merely copies
// cookbook directories to the right place.
package dir

import (
	"errors"
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/cookbook"
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

	files, err := cb.Files()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	c := rsync.Client{
		Archive:  true,
		Delete:   true,
		Compress: true,
		Verbose:  true,
	}
	return c.Copy(path.Join(dst, cb.Name), files...)
}

// Resolve copies the cookbook in the current directory to dst.
func (r Resolver) Resolve(dst string) error {
	return installCookbook(dst, ".")
}

// String returns the resolver's name.
func (r Resolver) String() string {
	return "Directory resolver"
}
