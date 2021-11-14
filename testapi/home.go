package testapi

import (
	"context"
	"net/http"
)

type Home struct {
	client      *Client
	computedURL string
}

func (h *Home) List(ctx context.Context) (message string, err error) {
	var resp struct {
		Message string `json:"message"`
	}

	if err := h.client.performJSONRequest(ctx, http.MethodGet, h.computedURL, nil, &resp); err != nil {
		return "", err
	}

	return resp.Message, nil
}

func (h *Home) Create(ctx context.Context, message string) error {
	var req struct {
		Message string `json:"message"`
	}

	return h.client.performJSONRequest(ctx, http.MethodPost, h.computedURL, req, nil)
}
