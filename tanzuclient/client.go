package tanzuclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client is a client for working with the TMC Web API.
// It is created by `NewClient`.
type Client struct {
	http           *http.Client
	baseURL        string
	token          AccessToken
	AcceptLanguage string
}

type AccessToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn string `json:"expires_in"`
}

func NewClient(url, apiToken *string) (*Client, error) {
	var token *AccessToken
	// Use apitoken (previously known as refresh token) to generate an access token.
	// Usually the access token is valid for a little less than 30minutes.
	loginURL := "https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize?refresh_token=" + *apiToken

	resp, err := http.Post(loginURL, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&token)

	return &Client{
		baseURL: *url,
		token:   *token,
		http: &http.Client{
			Timeout: time.Minute,
		},
	}, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.Token))

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}
