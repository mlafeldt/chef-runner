// Package util provides various utility functions.
package util

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// FileExist reports whether a file or directory exists.
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// BaseName - as the basename Unix tool - deletes any prefix ending with the
// last slash character present in a string, and a suffix, if given.
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

// InDir runs a function inside a specific directory.
func InDir(dir string, f func()) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := os.Chdir(dir); err != nil {
		panic(err)
	}

	f()

	if err := os.Chdir(wd); err != nil {
		panic(err)
	}
}

// InTestDir runs the passed function inside a temporary directory, which will
// be removed afterwards. Use it for isolated testing.
func InTestDir(f func()) {
	testDir, err := TempDir()
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(testDir)
	InDir(testDir, f)
}

// DownloadFile downloads a file from url and writes it to filename.
func DownloadFile(filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("HTTP error: " + resp.Status)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
