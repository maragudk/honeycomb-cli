package honeycomb

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// Dataset in Honeycomb.
type Dataset struct {
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Description     string `json:"description,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	LastWrittenAt   string `json:"last_written_at,omitempty"`
	RegularColumns  int    `json:"regular_columns,omitempty"`
	ExpandJSONDepth int    `json:"expand_json_depth,omitempty"`
}

// ListDatasets in the environment.
func (c *Client) ListDatasets(ctx context.Context) ([]Dataset, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/datasets", nil)
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

	var datasets []Dataset
	if err := json.NewDecoder(res.Body).Decode(&datasets); err != nil {
		return nil, err
	}
	return datasets, nil
}

// GetDataset by slug.
func (c *Client) GetDataset(ctx context.Context, slug string) (*Dataset, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/datasets/"+slug, nil)
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

	var dataset Dataset
	if err := json.NewDecoder(res.Body).Decode(&dataset); err != nil {
		return nil, err
	}
	return &dataset, nil
}

// CreateDatasetRequest for creating a dataset.
type CreateDatasetRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	ExpandJSONDepth int    `json:"expand_json_depth,omitempty"`
}

// CreateDataset with the given name.
func (c *Client) CreateDataset(ctx context.Context, create CreateDatasetRequest) (*Dataset, error) {
	body, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/1/datasets", bytes.NewReader(body))
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

	var dataset Dataset
	if err := json.NewDecoder(res.Body).Decode(&dataset); err != nil {
		return nil, err
	}
	return &dataset, nil
}

// DeleteDataset by slug.
func (c *Client) DeleteDataset(ctx context.Context, slug string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/1/datasets/"+slug, nil)
	if err != nil {
		return err
	}

	res, err := c.do(req)
	if err != nil {
		return err
	}
	_ = res.Body.Close()
	return nil
}
