package testapi

import (
	"context"
	"fmt"
	"net/http"
)

type Home struct {
	client      *Client
	computedURL string
}

func (h *Home) List(ctx context.Context) (messages []string, err error) {
	var resp struct {
		Messages []string `json:"messages"`
	}

	if err := h.client.performJSONRequest(ctx, http.MethodGet, h.computedURL, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

func (h *Home) Show(ctx context.Context, id string) (message string, err error) {
	var resp struct {
		Message string `json:"message"`
	}

	if err := h.client.performJSONRequest(ctx, http.MethodGet, fmt.Sprintf("%s/%s", h.computedURL, id), nil, &resp); err != nil {
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
