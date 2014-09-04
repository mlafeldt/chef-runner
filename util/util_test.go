package util_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
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

func TestDownloadFile(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("../testdata")))
	defer ts.Close()

	filename := "download.md"
	assert.NoError(t, util.DownloadFile(filename, ts.URL+"/README.md"))
	defer os.Remove(filename)

	data, _ := ioutil.ReadFile(filename)
	assert.Equal(t, "# Test Cookbook\n", string(data))
}
