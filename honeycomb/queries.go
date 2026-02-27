package honeycomb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// QuerySpec defines a Honeycomb query.
type QuerySpec struct {
	Calculations     []Calculation `json:"calculations,omitempty"`
	Breakdowns       []string      `json:"breakdowns,omitempty"`
	Filters          []Filter      `json:"filters,omitempty"`
	FilterCombination string       `json:"filter_combination,omitempty"`
	TimeRange        int           `json:"time_range,omitempty"`
	StartTime        int64         `json:"start_time,omitempty"`
	EndTime          int64         `json:"end_time,omitempty"`
	Orders           []Order       `json:"orders,omitempty"`
	Limit            int           `json:"limit,omitempty"`
}

// Calculation in a query (e.g. COUNT, AVG, P99).
type Calculation struct {
	Op     string `json:"op"`
	Column string `json:"column,omitempty"`
}

// Filter in a query.
type Filter struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  any    `json:"value"`
}

// Order for query results.
type Order struct {
	Op     string `json:"op,omitempty"`
	Column string `json:"column,omitempty"`
	Order  string `json:"order,omitempty"`
}

// QueryResponse from creating a query.
type QueryResponse struct {
	ID   string    `json:"id"`
	Spec QuerySpec `json:"query,omitempty"`
}

// QueryResultRequest to execute a query.
type QueryResultRequest struct {
	QueryID string `json:"query_id"`
}

// QueryResult from executing a query.
type QueryResult struct {
	ID       string `json:"id"`
	Complete bool   `json:"complete"`
	Data     struct {
		Results []map[string]any `json:"results"`
		Series  []any            `json:"series,omitempty"`
	} `json:"data"`
	Links struct {
		GraphURL string `json:"graph_image_url,omitempty"`
	} `json:"links,omitempty"`
}

// CreateQuery defines a query without executing it.
func (c *Client) CreateQuery(ctx context.Context, dataset string, spec QuerySpec) (*QueryResponse, error) {
	body, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/1/queries/"+dataset, bytes.NewReader(body))
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

	var query QueryResponse
	if err := json.NewDecoder(res.Body).Decode(&query); err != nil {
		return nil, err
	}
	return &query, nil
}

// CreateQueryResult executes a previously created query.
func (c *Client) CreateQueryResult(ctx context.Context, dataset, queryID string) (*QueryResult, error) {
	body, err := json.Marshal(QueryResultRequest{QueryID: queryID})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/1/query_results/"+dataset, bytes.NewReader(body))
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

	var result QueryResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetQueryResult polls for query results.
func (c *Client) GetQueryResult(ctx context.Context, dataset, resultID string) (*QueryResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v/1/query_results/%v/%v", c.baseURL, dataset, resultID), nil)
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

	var result QueryResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RunQuery creates a query, executes it, and polls until complete.
func (c *Client) RunQuery(ctx context.Context, dataset string, spec QuerySpec) (*QueryResult, error) {
	query, err := c.CreateQuery(ctx, dataset, spec)
	if err != nil {
		return nil, fmt.Errorf("creating query: %w", err)
	}

	result, err := c.CreateQueryResult(ctx, dataset, query.ID)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	for !result.Complete {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}

		result, err = c.GetQueryResult(ctx, dataset, result.ID)
		if err != nil {
			return nil, fmt.Errorf("polling query result: %w", err)
		}
	}

	return result, nil
}
