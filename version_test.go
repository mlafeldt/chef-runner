package main

import (
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestVersionString(t *testing.T) {
	GitVersion = ""
	assert.Equal(t, Version, VersionString())

	GitVersion = "some-git-version"
	assert.Equal(t, GitVersion, VersionString())
}
