package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestClient_ListMarkers(t *testing.T) {
	t.Run("returns markers for a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/markers/requests", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode([]honeycomb.Marker{
				{ID: "abc123", Type: "deploy", Message: "v1.0.0"},
				{ID: "def456", Type: "deploy", Message: "v1.1.0"},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		markers, err := c.ListMarkers(t.Context(), "requests")
		is.NotError(t, err)
		is.Equal(t, 2, len(markers))
		is.Equal(t, "v1.0.0", markers[0].Message)
	})
}

func TestClient_CreateMarker(t *testing.T) {
	t.Run("creates a marker for a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/markers/requests", r.URL.Path)
			is.Equal(t, http.MethodPost, r.Method)

			var req honeycomb.CreateMarkerRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			is.Equal(t, "deploy", req.Type)
			is.Equal(t, "v2.0.0 shipped", req.Message)

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(honeycomb.Marker{
				ID:      "new123",
				Type:    req.Type,
				Message: req.Message,
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		marker, err := c.CreateMarker(t.Context(), "requests", honeycomb.CreateMarkerRequest{
			Type:    "deploy",
			Message: "v2.0.0 shipped",
		})
		is.NotError(t, err)
		is.Equal(t, "new123", marker.ID)
		is.Equal(t, "deploy", marker.Type)
	})
}

func TestClient_DeleteMarker(t *testing.T) {
	t.Run("deletes a marker", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/markers/requests/abc123", r.URL.Path)
			is.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		err := c.DeleteMarker(t.Context(), "requests", "abc123")
		is.NotError(t, err)
	})
}
