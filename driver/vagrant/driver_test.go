package vagrant_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/driver/vagrant"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expect := "Vagrant driver (machine: some-machine)"
	actual := vagrant.NewDriver("some-machine").String()
	assert.Equal(t, expect, actual)
}

// TODO: test RunCommand machinery
