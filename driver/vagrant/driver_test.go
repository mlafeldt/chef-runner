package vagrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expect := "Vagrant driver (machine: some-machine)"
	actual := Driver{machine: "some-machine"}.String()
	assert.Equal(t, expect, actual)
}

// TODO: test driver machinery
