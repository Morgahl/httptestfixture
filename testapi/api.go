package testapi

import "net/http"

type TestAPI struct {
	Home *Home
}

func NewTestAPI(client *http.Client, url string) *TestAPI {
	internalClient := &Client{
		client: client,
	}
	return &TestAPI{
		&Home{client: internalClient, computedURL: url},
	}
}
