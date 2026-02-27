package honeycomb

import (
	"context"
	"encoding/json"
	"net/http"
)

// Column in a Honeycomb dataset.
type Column struct {
	ID          string `json:"id"`
	KeyName     string `json:"key_name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Hidden      bool   `json:"hidden,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	LastWritten string `json:"last_written,omitempty"`
}

// ListColumns for a dataset.
func (c *Client) ListColumns(ctx context.Context, dataset string) ([]Column, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/columns/"+dataset, nil)
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

	var columns []Column
	if err := json.NewDecoder(res.Body).Decode(&columns); err != nil {
		return nil, err
	}
	return columns, nil
}
