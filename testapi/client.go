package testapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
}

func (c *Client) performJSONRequest(ctx context.Context, method, path string, reqBody, respBody interface{}) error {
	req, err := makeJSONRequest(ctx, method, path, reqBody)
	if err != nil {
		return err
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if respBody != nil {
		req.Header.Set("Accept", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if err := handleJSONResponse(resp, &respBody); err != nil {
		return err
	}

	return nil
}

func makeJSONRequest(ctx context.Context, method, path string, reqBody interface{}) (*http.Request, error) {
	var body io.Reader
	if reqBody != nil {
		b := &bytes.Buffer{}
		if err := json.NewEncoder(b).Encode(reqBody); err != nil {
			return nil, err
		}
		body = b
	}

	return http.NewRequestWithContext(ctx, method, path, body)
}

func handleJSONResponse(resp *http.Response, respBody interface{}) error {
	defer resp.Body.Close()

	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
