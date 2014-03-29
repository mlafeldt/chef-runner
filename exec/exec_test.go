package exec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var history []string

func clearHistory() { history = []string{} }

func init() {
	cmdRunnerFunc = func(args []string) error {
		history = append(history, strings.Join(args, " "))
		return nil
	}
}

func TestExecuteCommand(t *testing.T) {
	defer clearHistory()

	ExecuteCommand([]string{"some", "test", "command"})

	if assert.Equal(t, 1, len(history)) {
		assert.Equal(t, "some test command", history[0])
	}
}
