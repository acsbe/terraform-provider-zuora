package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// NewRequest builds an authenticated HTTP request with bearer token.
func (c *Config) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}
	// build full URL
	url := fmt.Sprintf("%s%s", c.Endpoint, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
