// Package honeycomb provides a client for the Honeycomb.io API.
package honeycomb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client for the Honeycomb API.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// NewClient with the given API key and options.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: "https://api.honeycomb.io",
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Option for configuring a [Client].
type Option func(*Client)

// WithBaseURL sets the base URL for the API.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets the underlying HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.http = hc
	}
}

// APIError returned from the Honeycomb API.
type APIError struct {
	StatusCode int
	Status     int    `json:"status"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Err        string `json:"error"`
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("honeycomb API error (%d): %v", e.StatusCode, e.Detail)
	}
	if e.Err != "" {
		return fmt.Sprintf("honeycomb API error (%d): %v", e.StatusCode, e.Err)
	}
	return fmt.Sprintf("honeycomb API error (%d): %v", e.StatusCode, e.Title)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Honeycomb-Team", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		defer func() {
			_ = res.Body.Close()
		}()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("reading error response: %w", err)
		}

		apiErr := &APIError{StatusCode: res.StatusCode}
		if err := json.Unmarshal(body, apiErr); err != nil {
			return nil, fmt.Errorf("honeycomb API error (%d): %v", res.StatusCode, string(body))
		}
		return nil, apiErr
	}

	return res, nil
}
