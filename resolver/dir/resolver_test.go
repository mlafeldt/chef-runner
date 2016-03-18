package dir_test

import (
	"os"
	"path"
	"testing"

	"github.com/mlafeldt/chef-runner/resolver"
	. "github.com/mlafeldt/chef-runner/resolver/dir"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func TestResolverInterface(t *testing.T) {
	assert.Implements(t, (*resolver.Resolver)(nil), new(Resolver))
}

func TestResolve(t *testing.T) {
	defer util.TestChdir(t, "../../testdata")()

	cookbookPath := "test-cookbooks"
	defer os.RemoveAll(cookbookPath)

	assert.NoError(t, Resolver{}.Resolve(cookbookPath))

	expectFiles := []string{
		"practicingruby/README.md",
		"practicingruby/attributes",
		"practicingruby/metadata.rb",
		"practicingruby/recipes",
	}
	for _, f := range expectFiles {
		assert.True(t, util.FileExist(path.Join(cookbookPath, f)))
	}
}
