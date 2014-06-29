package metadata_test

import (
	"bytes"
	"testing"

	"github.com/mlafeldt/chef-runner.go/metadata"
	"github.com/stretchr/testify/assert"
)

func TestParse_Empty(t *testing.T) {
	metadata, err := metadata.Parse(bytes.NewBufferString(""))
	assert.NoError(t, err)
	if assert.NotNil(t, metadata) {
		assert.Equal(t, "", metadata.Name)
		assert.Equal(t, "", metadata.Version)
	}
}

func TestParse_Name(t *testing.T) {
	strings := []string{
		`name "cats"`,
		`name 'cats'`,
		` name   "cats" `,
	}
	for _, s := range strings {
		metadata, err := metadata.Parse(bytes.NewBufferString(s))
		assert.NoError(t, err)
		if assert.NotNil(t, metadata) {
			assert.Equal(t, "cats", metadata.Name)
		}
	}
}

func TestParse_Version(t *testing.T) {
	strings := []string{
		`version "1.2.3"`,
		`version '1.2.3'`,
		` version   "1.2.3" `,
	}
	for _, s := range strings {
		metadata, err := metadata.Parse(bytes.NewBufferString(s))
		assert.NoError(t, err)
		if assert.NotNil(t, metadata) {
			assert.Equal(t, "1.2.3", metadata.Version)
		}
	}
}

func TestParse_All(t *testing.T) {
	s := `
# some comment
name       "dogs"
maintainer "Pluto"
version    "2.0.0"
`
	metadata, err := metadata.Parse(bytes.NewBufferString(s))
	assert.NoError(t, err)
	if assert.NotNil(t, metadata) {
		assert.Equal(t, "dogs", metadata.Name)
		assert.Equal(t, "2.0.0", metadata.Version)
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
