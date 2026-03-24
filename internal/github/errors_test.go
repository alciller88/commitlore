package github

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	gh "github.com/google/go-github/v62/github"
	"github.com/stretchr/testify/assert"
)

func TestWrapAPIError_NotFound(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusNotFound},
	}
	err := wrapAPIError(ghErr, "alciller88", "repo-que-no-existe")
	assert.Contains(t, err.Error(), "repository not found: alciller88/repo-que-no-existe")
	assert.Contains(t, err.Error(), "check the name is correct")
	assert.NotContains(t, err.Error(), "api.github.com")
}

func TestWrapAPIError_Unauthorized(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusUnauthorized},
	}
	err := wrapAPIError(ghErr, "owner", "repo")
	assert.Contains(t, err.Error(), "authentication failed for owner/repo")
	assert.NotContains(t, err.Error(), "api.github.com")
}

func TestWrapAPIError_Forbidden(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
	}
	err := wrapAPIError(ghErr, "owner", "repo")
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestWrapAPIError_RateLimit(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Request:    &http.Request{URL: &url.URL{}},
		},
	}
	err := wrapAPIError(ghErr, "owner", "repo")
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.NotContains(t, err.Error(), "api.github.com")
}

func TestWrapAPIError_GenericError(t *testing.T) {
	err := wrapAPIError(fmt.Errorf("network timeout"), "owner", "repo")
	assert.Contains(t, err.Error(), "GitHub API error")
	assert.Contains(t, err.Error(), "network timeout")
}
