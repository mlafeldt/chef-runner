package runlist_test

import (
	"testing"

	"github.com/mlafeldt/chef-runner/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	. "github.com/mlafeldt/chef-runner/chef/runlist"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		cookbook  string
		recipes   []string
		runList   []string
		errString string
	}{
		{
			cookbook: "cats",
			recipes:  []string{},
			runList:  []string{},
		},
		{
			cookbook: "cats",
			recipes:  []string{"::foo"},
			runList:  []string{"cats::foo"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"recipes/foo.rb"},
			runList:  []string{"cats::foo"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"./recipes//foo.rb"},
			runList:  []string{"cats::foo"},
		},
		{
			cookbook: "",
			recipes:  []string{"dogs::bar"},
			runList:  []string{"dogs::bar"},
		},
		{
			cookbook: "",
			recipes:  []string{"dogs"},
			runList:  []string{"dogs"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"recipes/foo.rb", "::bar", "dogs::baz"},
			runList:  []string{"cats::foo", "cats::bar", "dogs::baz"},
		},
		{
			cookbook: "cats",
			recipes:  []string{"recipes/foo.rb,::bar,dogs::baz"},
			runList:  []string{"cats::foo", "cats::bar", "dogs::baz"},
		},
		// Check for errors
		{
			cookbook:  "",
			recipes:   []string{"::foo"},
			runList:   nil,
			errString: "cookbook name required",
		},
		{
			cookbook:  "",
			recipes:   []string{"recipes/foo.rb"},
			runList:   nil,
			errString: "cookbook name required",
		},
	}
	for _, test := range tests {
		runList, err := Build(test.recipes, test.cookbook)
		if test.errString == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errString)
		}
		assert.Equal(t, test.runList, runList)
	}
}
