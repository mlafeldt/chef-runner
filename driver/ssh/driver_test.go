package ssh_test

import (
	"testing"

	. "github.com/mlafeldt/chef-runner/driver/ssh"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	drv, err := NewDriver("some-user@some-host:1234")
	assert.NoError(t, err)
	assert.Equal(t, "SSH driver (host: some-host)", drv.String())
}

// TODO: test RunCommand machinery
