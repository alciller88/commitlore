package github

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	gh "github.com/google/go-github/v62/github"
)

// newTestClient creates a Client backed by a mock HTTP server.
func newTestClient(t *testing.T, mux *http.ServeMux) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(mux)
	baseURL, _ := url.Parse(server.URL + "/")

	ghClient := gh.NewClient(nil)
	ghClient.BaseURL = baseURL

	client := &Client{
		gh:    ghClient,
		owner: "testowner",
		repo:  "testrepo",
	}

	return client, server
}
