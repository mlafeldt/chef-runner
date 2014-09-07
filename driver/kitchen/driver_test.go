package kitchen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	expect := "Test Kitchen driver (instance: some-instance)"
	actual := Driver{instance: "some-instance"}.String()
	assert.Equal(t, expect, actual)
}

// TODO: test driver machinery
