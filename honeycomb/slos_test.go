package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestClient_ListSLOs(t *testing.T) {
	t.Run("returns SLOs for a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/slos/requests", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode([]honeycomb.SLO{
				{
					ID:               "slo1",
					Name:             "Request Latency",
					TargetPerMillion: 999000,
					TimePeriodDays:   30,
					SLI:              honeycomb.SLI{Alias: "sli.latency_ok"},
				},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		slos, err := c.ListSLOs(t.Context(), "requests")
		is.NotError(t, err)
		is.Equal(t, 1, len(slos))
		is.Equal(t, "Request Latency", slos[0].Name)
		is.Equal(t, 99.9, slos[0].TargetPercent())
	})
}

func TestClient_GetSLO(t *testing.T) {
	t.Run("returns a single SLO by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/slos/requests/slo1", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode(honeycomb.SLO{
				ID:               "slo1",
				Name:             "Request Latency",
				Description:      "P99 latency under 500ms",
				TargetPerMillion: 995000,
				TimePeriodDays:   7,
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		slo, err := c.GetSLO(t.Context(), "requests", "slo1")
		is.NotError(t, err)
		is.Equal(t, "Request Latency", slo.Name)
		is.Equal(t, "P99 latency under 500ms", slo.Description)
		is.Equal(t, 99.5, slo.TargetPercent())
	})
}
