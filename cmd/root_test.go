package cmd_test

import (
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/cmd"
)

func TestExecute(t *testing.T) {
	t.Run("returns 0 on success", func(t *testing.T) {
		code := cmd.Execute()
		is.Equal(t, 0, code)
	})
}
