package honeycomb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Marker in Honeycomb. Markers annotate points in time on graphs (e.g. deploys).
type Marker struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
}

// CreateMarkerRequest for creating a marker.
type CreateMarkerRequest struct {
	Type      string `json:"type,omitempty"`
	Message   string `json:"message,omitempty"`
	URL       string `json:"url,omitempty"`
	StartTime int64  `json:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty"`
}

// ListMarkers for a dataset. Use "__all__" for environment-wide markers.
func (c *Client) ListMarkers(ctx context.Context, dataset string) ([]Marker, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/markers/"+dataset, nil)
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

	var markers []Marker
	if err := json.NewDecoder(res.Body).Decode(&markers); err != nil {
		return nil, err
	}
	return markers, nil
}

// CreateMarker for a dataset. Use "__all__" for environment-wide markers.
func (c *Client) CreateMarker(ctx context.Context, dataset string, create CreateMarkerRequest) (*Marker, error) {
	body, err := json.Marshal(create)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/1/markers/"+dataset, bytes.NewReader(body))
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

	var marker Marker
	if err := json.NewDecoder(res.Body).Decode(&marker); err != nil {
		return nil, err
	}
	return &marker, nil
}

// DeleteMarker by ID for a dataset.
func (c *Client) DeleteMarker(ctx context.Context, dataset, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%v/1/markers/%v/%v", c.baseURL, dataset, id), nil)
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
