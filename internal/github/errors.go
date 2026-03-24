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

func wrapAPIError(err error, owner, repo string) error {
	ref := owner + "/" + repo

	var ghErr *gh.ErrorResponse
	if errors.As(err, &ghErr) {
		switch ghErr.Response.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("repository not found: %s\n(check the name is correct or set GITHUB_TOKEN for private repos)", ref)
		case http.StatusUnauthorized, http.StatusForbidden:
			return fmt.Errorf("GitHub API authentication failed for %s\n(check that GITHUB_TOKEN is valid)", ref)
		case http.StatusTooManyRequests:
			return fmt.Errorf("GitHub API rate limit exceeded\n(authenticate with GITHUB_TOKEN for higher limits)")
		}
	}
	return fmt.Errorf("GitHub API error: %w", err)
}
