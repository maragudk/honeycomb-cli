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

func TestColumnsListCommand(t *testing.T) {
	t.Run("lists columns in a table", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/columns/requests", r.URL.Path)
			_ = json.NewEncoder(w).Encode([]honeycomb.Column{
				{ID: "1", KeyName: "duration_ms", Type: "float", Description: "How long it took"},
				{ID: "2", KeyName: "status_code", Type: "integer"},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"columns", "list", "--dataset", "requests", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "duration_ms"))
		is.True(t, contains(output, "status_code"))
		is.True(t, contains(output, "KEY NAME"))
	})

	t.Run("lists columns as JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode([]honeycomb.Column{
				{ID: "1", KeyName: "duration_ms", Type: "float"},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"columns", "list", "--dataset", "requests", "--json", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		var columns []honeycomb.Column
		is.NotError(t, json.Unmarshal(buf.Bytes(), &columns))
		is.Equal(t, 1, len(columns))
		is.Equal(t, "duration_ms", columns[0].KeyName)
	})
}
