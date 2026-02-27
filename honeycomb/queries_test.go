package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestClient_CreateQuery(t *testing.T) {
	t.Run("creates a query definition", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/queries/requests", r.URL.Path)
			is.Equal(t, http.MethodPost, r.Method)

			var spec honeycomb.QuerySpec
			_ = json.NewDecoder(r.Body).Decode(&spec)
			is.Equal(t, 1, len(spec.Calculations))
			is.Equal(t, "COUNT", spec.Calculations[0].Op)

			_ = json.NewEncoder(w).Encode(honeycomb.QueryResponse{
				ID: "query-abc",
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		query, err := c.CreateQuery(t.Context(), "requests", honeycomb.QuerySpec{
			Calculations: []honeycomb.Calculation{{Op: "COUNT"}},
			TimeRange:    7200,
		})
		is.NotError(t, err)
		is.Equal(t, "query-abc", query.ID)
	})
}

func TestClient_RunQuery(t *testing.T) {
	t.Run("creates, executes, and polls for query results", func(t *testing.T) {
		var pollCount atomic.Int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/1/queries/requests":
				_ = json.NewEncoder(w).Encode(honeycomb.QueryResponse{ID: "q1"})

			case r.Method == http.MethodPost && r.URL.Path == "/1/query_results/requests":
				var req honeycomb.QueryResultRequest
				_ = json.NewDecoder(r.Body).Decode(&req)
				is.Equal(t, "q1", req.QueryID)

				_ = json.NewEncoder(w).Encode(honeycomb.QueryResult{
					ID:       "r1",
					Complete: false,
				})

			case r.Method == http.MethodGet && r.URL.Path == "/1/query_results/requests/r1":
				count := pollCount.Add(1)
				complete := count >= 2

				result := honeycomb.QueryResult{
					ID:       "r1",
					Complete: complete,
				}
				if complete {
					result.Data.Results = []map[string]any{
						{"COUNT": float64(42)},
					}
				}
				_ = json.NewEncoder(w).Encode(result)
			}
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		result, err := c.RunQuery(t.Context(), "requests", honeycomb.QuerySpec{
			Calculations: []honeycomb.Calculation{{Op: "COUNT"}},
			TimeRange:    3600,
		})
		is.NotError(t, err)
		is.True(t, result.Complete)
		is.Equal(t, 1, len(result.Data.Results))
		is.Equal(t, float64(42), result.Data.Results[0]["COUNT"].(float64))
	})
}
