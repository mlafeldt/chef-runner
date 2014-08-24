package librarian_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/resolver/librarian"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	expect := []string{"librarian-chef", "install", "--path", "a/b/c"}
	actual := librarian.Command("a/b/c")
	assert.Equal(t, expect, actual)
}
