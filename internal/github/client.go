package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	gh "github.com/google/go-github/v62/github"
)

// ownerRepoPattern matches "owner/repo" format (no slashes in components).
var ownerRepoPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)

// githubURLPattern extracts owner/repo from GitHub URLs.
var githubURLPattern = regexp.MustCompile(
	`^https?://github\.com/([a-zA-Z0-9._-]+)/([a-zA-Z0-9._-]+?)(?:\.git)?/?$`,
)

// Client wraps the go-github client for read-only operations.
type Client struct {
	gh    *gh.Client
	owner string
	repo  string
}

// NewClient creates a GitHub API client for the given owner/repo.
// If GITHUB_TOKEN is set, it authenticates for private repo access.
func NewClient(owner, repo string) *Client {
	var httpClient *http.Client

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		httpClient = &http.Client{
			Transport: &tokenTransport{token: token},
		}
	}

	return &Client{
		gh:    gh.NewClient(httpClient),
		owner: owner,
		repo:  repo,
	}
}

// tokenTransport adds Bearer token authentication to HTTP requests.
type tokenTransport struct {
	token string
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(context.Background())
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}

// IsRemoteRepo returns true if the value looks like a GitHub reference
// (owner/repo or https://github.com/owner/repo URL).
func IsRemoteRepo(value string) bool {
	return ownerRepoPattern.MatchString(value) || githubURLPattern.MatchString(value)
}

// ParseRepoRef extracts owner and repo from an owner/repo string or GitHub URL.
func ParseRepoRef(value string) (owner, repo string, err error) {
	if matches := githubURLPattern.FindStringSubmatch(value); len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	if ownerRepoPattern.MatchString(value) {
		parts := strings.SplitN(value, "/", 2)
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("invalid GitHub reference %q: use owner/repo or https://github.com/owner/repo", value)
}
