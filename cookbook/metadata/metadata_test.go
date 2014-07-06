package metadata_test

import (
	"bytes"
	"testing"

	"github.com/mlafeldt/chef-runner.go/cookbook/metadata"
	"github.com/stretchr/testify/assert"
)

var parseTests = []struct {
	in            string
	name, version string
}{
	{"", "", ""},
	{`name "cats"`, "cats", ""},
	{`name 'cats'`, "cats", ""},
	{` name   "cats" `, "cats", ""},
	{`version "1.2.3"`, "", "1.2.3"},
	{`version '1.2.3'`, "", "1.2.3"},
	{` version   "1.2.3" `, "", "1.2.3"},
	{`
# some comment
name       "dogs"
maintainer "Pluto"
version    "2.0.0"`, "dogs", "2.0.0"},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		metadata, err := metadata.Parse(bytes.NewBufferString(test.in))
		assert.NoError(t, err)
		if assert.NotNil(t, metadata) {
			assert.Equal(t, test.name, metadata.Name)
			assert.Equal(t, test.version, metadata.Version)
		}
	}
}

func TestParseFile(t *testing.T) {
	metadata, err := metadata.ParseFile("testdata/metadata.rb")
	assert.NoError(t, err)
	if assert.NotNil(t, metadata) {
		assert.Equal(t, "practicingruby", metadata.Name)
		assert.Equal(t, "1.3.1", metadata.Version)
	}
}
