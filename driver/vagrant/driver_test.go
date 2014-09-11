package vagrant

import (
	"testing"

	"github.com/mlafeldt/chef-runner/driver"
	"github.com/stretchr/testify/assert"
)

func TestDriverInterface(t *testing.T) {
	assert.Implements(t, (*driver.Driver)(nil), new(Driver))
}

func TestString(t *testing.T) {
	expect := "Vagrant driver (machine: some-machine)"
	actual := Driver{machine: "some-machine"}.String()
	assert.Equal(t, expect, actual)
}

// TODO: test driver machinery
