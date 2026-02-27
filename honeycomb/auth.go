package honeycomb

import (
	"context"
	"encoding/json"
	"net/http"
)

// AuthResponse from the Honeycomb auth endpoint.
type AuthResponse struct {
	Team         AuthTeam        `json:"team"`
	Environment  AuthEnvironment `json:"environment"`
	APIKeyAccess map[string]bool `json:"api_key_access"`
}

// AuthTeam information.
type AuthTeam struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// AuthEnvironment information.
type AuthEnvironment struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Auth verifies the API key and returns team and environment information.
func (c *Client) Auth(ctx context.Context) (*AuthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/1/auth", nil)
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

	var auth AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&auth); err != nil {
		return nil, err
	}
	return &auth, nil
}
