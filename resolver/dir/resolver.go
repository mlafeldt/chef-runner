// Package dir implements a cookbook dependency resolver based on rsync. This
// resolver is pretty basic in that it only copies cookbook directories.
package dir

import (
	"errors"
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/cookbook"
	"github.com/mlafeldt/chef-runner/rsync"
)

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
func Resolve(dst string) error {
	return installCookbook(dst, ".")
}
