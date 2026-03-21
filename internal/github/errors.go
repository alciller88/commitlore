package github

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	gh "github.com/google/go-github/v62/github"
)

// ErrTokenRequired indicates that a GITHUB_TOKEN is needed for the operation.
var ErrTokenRequired = errors.New("GITHUB_TOKEN environment variable is required for this operation")

// HasToken returns true if a GitHub token is configured.
func HasToken() bool {
	return os.Getenv("GITHUB_TOKEN") != ""
}

func wrapAPIError(err error) error {
	var ghErr *gh.ErrorResponse
	if errors.As(err, &ghErr) {
		switch ghErr.Response.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("repository not found (check owner/repo or set GITHUB_TOKEN for private repos): %w", err)
		case http.StatusUnauthorized, http.StatusForbidden:
			return fmt.Errorf("GitHub API authentication failed (check GITHUB_TOKEN): %w", err)
		case http.StatusTooManyRequests:
			return fmt.Errorf("GitHub API rate limit exceeded (authenticate with GITHUB_TOKEN for higher limits): %w", err)
		}
	}
	return fmt.Errorf("GitHub API error: %w", err)
}
