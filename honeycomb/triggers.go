package honeycomb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Trigger (alert) in Honeycomb.
type Trigger struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Disabled    bool             `json:"disabled"`
	Frequency   int              `json:"frequency,omitempty"`
	Threshold   TriggerThreshold `json:"threshold"`
	CreatedAt   string           `json:"created_at,omitempty"`
	UpdatedAt   string           `json:"updated_at,omitempty"`
}

// TriggerThreshold defines when a trigger fires.
type TriggerThreshold struct {
	Op    string  `json:"op"`
	Value float64 `json:"value"`
}

// ListTriggers for a dataset.
func (c *Client) ListTriggers(ctx context.Context, dataset string) ([]Trigger, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/triggers/"+dataset, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var triggers []Trigger
	if err := json.NewDecoder(res.Body).Decode(&triggers); err != nil {
		return nil, err
	}
	return triggers, nil
}

// GetTrigger by ID for a dataset.
func (c *Client) GetTrigger(ctx context.Context, dataset, id string) (*Trigger, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v/1/triggers/%v/%v", c.baseURL, dataset, id), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var trigger Trigger
	if err := json.NewDecoder(res.Body).Decode(&trigger); err != nil {
		return nil, err
	}
	return &trigger, nil
}
