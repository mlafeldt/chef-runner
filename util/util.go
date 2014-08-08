// Package util provides various utility functions.
package util

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// FileExist reports whether a file or directory exists.
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// BaseName, as the basename Unix tool, deletes any prefix ending with the last
// slash character present in a string, and a suffix, if given.
func BaseName(s, suffix string) string {
	base := path.Base(s)
	if suffix != "" {
		base = strings.TrimSuffix(base, suffix)
	}
	return base
}

// TempDir creates a new temporary directory to be used by chef-runner.
func TempDir() (string, error) {
	return ioutil.TempDir("", "chef-runner-")
}
