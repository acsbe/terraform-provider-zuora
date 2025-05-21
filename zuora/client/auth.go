package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Config struct {
	ClientID     string
	ClientSecret string
	Endpoint     string
	HTTPClient   *http.Client

	token       string
	tokenExpiry time.Time
}

// getToken returns a cached token or fetches a new one if expired.
func (c *Config) getToken(ctx context.Context) (string, error) {
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	data := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"client_credentials"},
	}

	req, err := http.NewRequestWithContext(ctx,
		"POST",
		fmt.Sprintf("%s/oauth/token", c.Endpoint),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	c.token = tr.AccessToken
	// subtract a small buffer so we never hit expiry exactly
	c.tokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn-10) * time.Second)
	return c.token, nil
}
