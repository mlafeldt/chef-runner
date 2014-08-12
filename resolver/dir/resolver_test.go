package dir_test

import (
	"os"
	"path"
	"testing"

	"github.com/mlafeldt/chef-runner/resolver/dir"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

const CookbookPath = "test-cookbooks"

func TestResolve(t *testing.T) {
	if err := os.Chdir("../../testdata"); err != nil {
		panic(err)
	}

	defer os.RemoveAll(CookbookPath)

	assert.NoError(t, dir.Resolver{}.Resolve(CookbookPath))

	expectFiles := []string{
		"practicingruby/README.md",
		"practicingruby/attributes",
		"practicingruby/metadata.rb",
		"practicingruby/recipes",
	}
	for _, f := range expectFiles {
		assert.True(t, util.FileExist(path.Join(CookbookPath, f)))
	}
}
