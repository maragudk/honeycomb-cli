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

func TestSLOsListCommand(t *testing.T) {
	t.Run("lists SLOs in a table", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/slos/requests", r.URL.Path)
			_ = json.NewEncoder(w).Encode([]honeycomb.SLO{
				{ID: "slo1", Name: "Latency", TargetPerMillion: 999000, TimePeriodDays: 30},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"slos", "list", "--dataset", "requests", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "Latency"))
		is.True(t, contains(output, "99.90%"))
		is.True(t, contains(output, "30d"))
	})
}

func TestSLOsGetCommand(t *testing.T) {
	t.Run("shows SLO details", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/slos/requests/slo1", r.URL.Path)
			_ = json.NewEncoder(w).Encode(honeycomb.SLO{
				ID:               "slo1",
				Name:             "Latency",
				Description:      "P99 under 500ms",
				TargetPerMillion: 999000,
				TimePeriodDays:   30,
				SLI:              honeycomb.SLI{Alias: "sli.latency"},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"slos", "get", "slo1", "--dataset", "requests", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "Latency"))
		is.True(t, contains(output, "99.90%"))
		is.True(t, contains(output, "sli.latency"))
	})
}
