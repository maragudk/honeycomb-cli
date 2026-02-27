package honeycomb_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"maragu.dev/is"

	"github.com/maragudk/honeycomb-cli/honeycomb"
)

func TestClient_ListTriggers(t *testing.T) {
	t.Run("returns triggers for a dataset", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/triggers/requests", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode([]honeycomb.Trigger{
				{
					ID:   "t1",
					Name: "High Error Rate",
					Threshold: honeycomb.TriggerThreshold{
						Op:    ">",
						Value: 0.05,
					},
					Frequency: 300,
				},
				{
					ID:       "t2",
					Name:     "Old Alert",
					Disabled: true,
				},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		triggers, err := c.ListTriggers(t.Context(), "requests")
		is.NotError(t, err)
		is.Equal(t, 2, len(triggers))
		is.Equal(t, "High Error Rate", triggers[0].Name)
		is.Equal(t, ">", triggers[0].Threshold.Op)
		is.True(t, triggers[1].Disabled)
	})
}

func TestClient_GetTrigger(t *testing.T) {
	t.Run("returns a single trigger by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			is.Equal(t, "/1/triggers/requests/t1", r.URL.Path)
			is.Equal(t, http.MethodGet, r.Method)

			_ = json.NewEncoder(w).Encode(honeycomb.Trigger{
				ID:          "t1",
				Name:        "High Error Rate",
				Description: "Fires when error rate exceeds 5%",
				Frequency:   300,
				Threshold: honeycomb.TriggerThreshold{
					Op:    ">",
					Value: 0.05,
				},
			})
		}))
		defer server.Close()

		c := honeycomb.NewClient("test-key", honeycomb.WithBaseURL(server.URL))
		trigger, err := c.GetTrigger(t.Context(), "requests", "t1")
		is.NotError(t, err)
		is.Equal(t, "High Error Rate", trigger.Name)
		is.Equal(t, "Fires when error rate exceeds 5%", trigger.Description)
		is.Equal(t, 300, trigger.Frequency)
	})
}
