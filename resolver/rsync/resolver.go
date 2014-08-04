package rsync

import (
	"os"
	"path"

	"github.com/mlafeldt/chef-runner/cookbook"
	"github.com/mlafeldt/chef-runner/rsync"
)

func InstallCookbook(dst, src string) error {
	cb, err := cookbook.NewCookbook(src)
	if err != nil {
		return err
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
