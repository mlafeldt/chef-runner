package dir_test

import (
	"os"
	"path"
	"testing"

	. "github.com/mlafeldt/chef-runner/resolver/dir"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	util.InDir("../../testdata", func() {
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
	})
}
