package cmd_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/cmd"
	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestAuthCommand(t *testing.T) {
	t.Run("displays team and environment info", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(honeycomb.AuthResponse{
				Team: honeycomb.AuthTeam{
					Name: "Acme Corp",
				},
				Environment: honeycomb.AuthEnvironment{
					Name: "Production",
				},
				APIKeyAccess: map[string]bool{
					"events": true,
				},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"auth", "--api-key", "test-key", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "Acme Corp"))
		is.True(t, contains(output, "Production"))
		is.True(t, contains(output, "events"))
	})

	t.Run("returns error when API key is missing", func(t *testing.T) {
		root := cmd.NewRootCommand()
		root.SetArgs([]string{"auth"})

		// Clear the env var for this test
		t.Setenv("HONEYCOMB_API_KEY", "")

		err := root.Execute()
		is.True(t, err != nil)
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
