package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestNewClient(t *testing.T) {
	t.Run("uses default base URL", func(t *testing.T) {
		c := honeycomb.NewClient("test-key")
		is.NotNil(t, c)
	})

	t.Run("accepts custom base URL", func(t *testing.T) {
		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL("https://custom.example.com"))
		is.NotNil(t, c)
	})
}

func TestClient_Auth(t *testing.T) {
	t.Run("returns auth info for a valid API key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/auth", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)
			is.Equal(t, "test-key", r.Header.Get("X-Honeycomb-Team"))

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(honeycomb.AuthResponse{
				Team: honeycomb.AuthTeam{
					Name: "My Team",
					Slug: "my-team",
				},
				Environment: honeycomb.AuthEnvironment{
					Name: "Production",
					Slug: "production",
				},
				APIKeyAccess: map[string]bool{
					"events":  true,
					"markers": true,
				},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		auth, err := c.Auth(t.Context())
		is.NotError(t, err)
		is.Equal(t, "My Team", auth.Team.Name)
		is.Equal(t, "my-team", auth.Team.Slug)
		is.Equal(t, "Production", auth.Environment.Name)
		is.Equal(t, "production", auth.Environment.Slug)
		is.True(t, auth.APIKeyAccess["events"])
		is.True(t, auth.APIKeyAccess["markers"])
	})

	t.Run("returns error for invalid API key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": 401,
				"error":  "unknown API key - check your credentials",
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("bad-key", honeycomb.WithBaseURL(server.URL))
		_, err := c.Auth(t.Context())
		is.True(t, err != nil)
	})
}
