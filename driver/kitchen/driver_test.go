package kitchen

import (
	"testing"

	"github.com/mlafeldt/chef-runner/driver"
	"github.com/stretchr/testify/assert"
)

func TestDriverInterface(t *testing.T) {
	assert.Implements(t, (*driver.Driver)(nil), new(Driver))
}

func TestString(t *testing.T) {
	expect := "Test Kitchen driver (instance: some-instance)"
	actual := Driver{instance: "some-instance"}.String()
	assert.Equal(t, expect, actual)
}

// TODO: test driver machinery
