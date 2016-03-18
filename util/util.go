// Package util provides various utility functions.
package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// FileExist reports whether a file or directory exists.
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// BaseName - as the basename Unix tool - deletes any prefix ending with the
// last slash character present in a string, and a suffix, if given.
func BaseName(s, suffix string) string {
	base := filepath.Base(s)
	if suffix != "" {
		base = strings.TrimSuffix(base, suffix)
	}
	return base
}

func TestChdir(t *testing.T, dir string) func() {
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("error: %s", err)
	}
	return func() { os.Chdir(old) }
}

func TestTempDir(t *testing.T) func() {
	tmp, err := ioutil.TempDir("", "chef-runner-")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	f := TestChdir(t, tmp)
	return func() { f(); os.RemoveAll(tmp) }
}
