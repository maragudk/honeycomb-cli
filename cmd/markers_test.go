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

func TestMarkersListCommand(t *testing.T) {
	t.Run("lists markers in a table", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/markers/__all__", r.URL.Path)
			_ = json.NewEncoder(w).Encode([]honeycomb.Marker{
				{ID: "abc", Type: "deploy", Message: "v1.0.0"},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"markers", "list", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "deploy"))
		is.True(t, contains(output, "v1.0.0"))
	})

	t.Run("lists markers for a specific dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/markers/my-dataset", r.URL.Path)
			_ = json.NewEncoder(w).Encode([]honeycomb.Marker{})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"markers", "list", "--dataset", "my-dataset", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)
	})
}

func TestMarkersCreateCommand(t *testing.T) {
	t.Run("creates a deploy marker", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, http.MethodPost, r.Method)

			var req honeycomb.CreateMarkerRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			is.Equal(t, "deploy", req.Type)
			is.Equal(t, "v2.0.0", req.Message)

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(honeycomb.Marker{
				ID:   "new123",
				Type: req.Type,
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"markers", "create", "--type", "deploy", "--message", "v2.0.0", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		is.True(t, contains(buf.String(), "new123"))
	})
}

func TestMarkersDeleteCommand(t *testing.T) {
	t.Run("deletes a marker", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, http.MethodDelete, r.Method)
			is.Equal(t, "/1/markers/__all__/abc123", r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"markers", "delete", "abc123", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		is.True(t, contains(buf.String(), "abc123"))
	})
}
