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

func TestTriggersListCommand(t *testing.T) {
	t.Run("lists triggers in a table", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/triggers/requests", r.URL.Path)
			_ = json.NewEncoder(w).Encode([]honeycomb.Trigger{
				{
					ID:   "t1",
					Name: "High Error Rate",
					Threshold: honeycomb.TriggerThreshold{
						Op:    ">",
						Value: 0.05,
					},
				},
				{
					ID:       "t2",
					Name:     "Disabled Alert",
					Disabled: true,
				},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"triggers", "list", "--dataset", "requests", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "High Error Rate"))
		is.True(t, contains(output, "Disabled Alert"))
		is.True(t, contains(output, "yes"))
	})
}

func TestTriggersGetCommand(t *testing.T) {
	t.Run("shows trigger details", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/triggers/requests/t1", r.URL.Path)
			_ = json.NewEncoder(w).Encode(honeycomb.Trigger{
				ID:          "t1",
				Name:        "High Error Rate",
				Description: "Fires at 5% errors",
				Frequency:   300,
				Threshold: honeycomb.TriggerThreshold{
					Op:    ">",
					Value: 0.05,
				},
			})
		}))
		defer server.Close()

		var buf bytes.Buffer
		root := cmd.NewRootCommand()
		root.SetOut(&buf)
		root.SetArgs([]string{"triggers", "get", "t1", "--dataset", "requests", "--api-key", "test", "--api-url", server.URL})

		err := root.Execute()
		is.NotError(t, err)

		output := buf.String()
		is.True(t, contains(output, "High Error Rate"))
		is.True(t, contains(output, "Fires at 5% errors"))
		is.True(t, contains(output, "300"))
	})
}
