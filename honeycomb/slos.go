package honeycomb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SLO in Honeycomb.
type SLO struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	SLI             SLI     `json:"sli"`
	TargetPerMillion int    `json:"target_per_million"`
	TimePeriodDays  int     `json:"time_period_days"`
	CreatedAt       string  `json:"created_at,omitempty"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
}

// TargetPercent returns the target as a human-readable percentage.
func (s *SLO) TargetPercent() float64 {
	return float64(s.TargetPerMillion) / 10000
}

// SLI (Service Level Indicator) definition.
type SLI struct {
	Alias string `json:"alias"`
}

// ListSLOs for a dataset.
func (c *Client) ListSLOs(ctx context.Context, dataset string) ([]SLO, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/slos/"+dataset, nil)
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

	var slos []SLO
	if err := json.NewDecoder(res.Body).Decode(&slos); err != nil {
		return nil, err
	}
	return slos, nil
}

// GetSLO by ID for a dataset.
func (c *Client) GetSLO(ctx context.Context, dataset, id string) (*SLO, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v/1/slos/%v/%v", c.baseURL, dataset, id), nil)
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

	var slo SLO
	if err := json.NewDecoder(res.Body).Decode(&slo); err != nil {
		return nil, err
	}
	return &slo, nil
}
