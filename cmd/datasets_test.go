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

func TestDatasetsListCommand(t *testing.T) {
	t.Run("lists datasets in a table", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode([]honeycomb.Dataset{
				{Name: "requests", Slug: "requests", Description: "HTTP requests"},
				{Name: "errors", Slug: "errors"},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"datasets", "list", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "requests"))
		is.True(t, contains(output, "errors"))
		is.True(t, contains(output, "NAME"))
	})

	t.Run("lists datasets as JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode([]honeycomb.Dataset{
				{Name: "requests", Slug: "requests"},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"datasets", "list", "--json", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		var datasets []honeycomb.Dataset
		is.NotError(t, json.Unmarshal(buf.Bytes(), &datasets))
		is.Equal(t, 1, len(datasets))
		is.Equal(t, "requests", datasets[0].Name)
	})
}

func TestDatasetsGetCommand(t *testing.T) {
	t.Run("shows dataset details", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(honeycomb.Dataset{
				Name:        "requests",
				Slug:        "requests",
				Description: "HTTP requests",
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"datasets", "get", "requests", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "requests"))
		is.True(t, contains(output, "HTTP requests"))
	})
}

func TestDatasetsCreateCommand(t *testing.T) {
	t.Run("creates a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(honeycomb.Dataset{
				Name: "new-dataset",
				Slug: "new-dataset",
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"datasets", "create", "--name", "new-dataset", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		is.True(t, contains(buf.String(), "new-dataset"))
	})
}

func TestDatasetsDeleteCommand(t *testing.T) {
	t.Run("deletes a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"datasets", "delete", "old-dataset", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		is.True(t, contains(buf.String(), "old-dataset"))
	})
}
