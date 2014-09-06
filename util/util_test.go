package util_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func TestFileExist(t *testing.T) {
	filename := "some-file"
	assert.False(t, util.FileExist(filename))

	f, _ := os.Create(filename)
	f.Close()
	defer os.Remove(filename)
	assert.True(t, util.FileExist(filename))
}

func TestBaseName(t *testing.T) {
	tests := []struct {
		in     string
		suffix string
		out    string
	}{
		{"", "", "."},
		{"a", "", "a"},
		{"a/b", "", "b"},
		{"/a/b/c", "", "c"},
		{"a.x", ".x", "a"},
		{"a/b.x", ".x", "b"},
		{"a/b.x", ".y", "b.x"},
	}
	for _, test := range tests {
		assert.Equal(t, test.out, util.BaseName(test.in, test.suffix))
	}
}

func TestTempDir(t *testing.T) {
	dir, err := util.TempDir()
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	assert.True(t, strings.HasPrefix(path.Base(dir), "chef-runner-"))

	m, err := os.Stat(dir)
	assert.NoError(t, err)
	assert.True(t, m.IsDir())
}

func TestInDir(t *testing.T) {
	var wd1, wd2, wd3 string
	wd1, _ = os.Getwd()

	util.InDir("../testdata", func() {
		wd2, _ = os.Getwd()
	})

	abs, _ := filepath.Abs("../testdata")
	assert.Equal(t, abs, wd2)

	wd3, _ = os.Getwd()
	assert.Equal(t, wd3, wd1)
}

func TestInTestDir(t *testing.T) {
	wd, _ := os.Getwd()
	var testDir string

	util.InTestDir(func() {
		testDir, _ = os.Getwd()
		assert.NotEqual(t, testDir, wd)
		assert.NoError(t, ioutil.WriteFile("some-test-file", []byte{}, 0644))
	})

	wd2, _ := os.Getwd()
	assert.Equal(t, wd, wd2)
	assert.False(t, util.FileExist(testDir))
}

func TestDownloadFile(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("../testdata")))
	defer ts.Close()

	filename := "download.md"
	assert.NoError(t, util.DownloadFile(filename, ts.URL+"/README.md"))
	defer os.Remove(filename)

	data, _ := ioutil.ReadFile(filename)
	assert.Equal(t, "# Test Cookbook\n", string(data))
}
