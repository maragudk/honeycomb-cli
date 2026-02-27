package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestClient_ListColumns(t *testing.T) {
	t.Run("returns columns for a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/columns/requests", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode([]honeycomb.Column{
				{ID: "1", KeyName: "duration_ms", Type: "float", Description: "Request duration"},
				{ID: "2", KeyName: "status_code", Type: "integer"},
				{ID: "3", KeyName: "trace.trace_id", Type: "string", Hidden: true},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		columns, err := c.ListColumns(t.Context(), "requests")
		is.NotError(t, err)
		is.Equal(t, 3, len(columns))
		is.Equal(t, "duration_ms", columns[0].KeyName)
		is.Equal(t, "float", columns[0].Type)
		is.Equal(t, "Request duration", columns[0].Description)
		is.True(t, columns[2].Hidden)
	})
}
