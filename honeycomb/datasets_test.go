package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestClient_ListDatasets(t *testing.T) {
	t.Run("returns list of datasets", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/datasets", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode([]honeycomb.Dataset{
				{Name: "requests", Slug: "requests"},
				{Name: "errors", Slug: "errors"},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		datasets, err := c.ListDatasets(t.Context())
		is.NotError(t, err)
		is.Equal(t, 2, len(datasets))
		is.Equal(t, "requests", datasets[0].Name)
		is.Equal(t, "errors", datasets[1].Name)
	})
}

func TestClient_GetDataset(t *testing.T) {
	t.Run("returns dataset by slug", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/datasets/requests", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode(honeycomb.Dataset{
				Name:        "requests",
				Slug:        "requests",
				Description: "HTTP requests",
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		dataset, err := c.GetDataset(t.Context(), "requests")
		is.NotError(t, err)
		is.Equal(t, "requests", dataset.Name)
		is.Equal(t, "HTTP requests", dataset.Description)
	})
}

func TestClient_CreateDataset(t *testing.T) {
	t.Run("creates a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/datasets", r.URL.Path)
			is.Equal(t, http.MethodPost, r.Method)

			var req honeycomb.CreateDatasetRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			is.Equal(t, "new-dataset", req.Name)

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(honeycomb.Dataset{
				Name: "new-dataset",
				Slug: "new-dataset",
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		dataset, err := c.CreateDataset(t.Context(), honeycomb.CreateDatasetRequest{Name: "new-dataset"})
		is.NotError(t, err)
		is.Equal(t, "new-dataset", dataset.Name)
	})
}

func TestClient_DeleteDataset(t *testing.T) {
	t.Run("deletes a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/datasets/old-dataset", r.URL.Path)
			is.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		err := c.DeleteDataset(t.Context(), "old-dataset")
		is.NotError(t, err)
	})
}
