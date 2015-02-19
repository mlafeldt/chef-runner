// Package util provides various utility functions.
package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

// WriteTimestampFile writes the current time as Unix timestamp to a file.
func WriteTimestampFile(filename string) error {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return ioutil.WriteFile(filename, []byte(ts), 0644)
}

// ReadTimestampFile reads the Unix timestamp from a file created with
// WriteTimestampFile.
func ReadTimestampFile(filename string) (int64, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(data), 10, 64)
}

// FileModTime returns a file's modification time as Unix timestamp.
func FileModTime(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.ModTime().Unix(), nil
}
