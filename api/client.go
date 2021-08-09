package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client is a client for working with the TMC Web API.
// It is created by `NewClient` and `Authenticator.NewClient`.
type Client struct {
	http    *http.Client
	baseURL string

	AutoRetry      bool
	AcceptLanguage string
}

func (c *Client) get(url string, result interface{}) error {
	for {
		req, err := http.NewRequest("GET", url, nil)
		if c.AcceptLanguage != "" {
			req.Header.Set("Accept-Language", c.AcceptLanguage)
		}
		if err != nil {
			return err
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == rateLimitExceededStatusCode && c.AutoRetry {
			time.Sleep(retryDuration(resp))
			continue
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			return c.decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}
