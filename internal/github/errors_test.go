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
	err := wrapAPIError(ghErr)
	assert.Contains(t, err.Error(), "repository not found")
}

func TestWrapAPIError_Unauthorized(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusUnauthorized},
	}
	err := wrapAPIError(ghErr)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestWrapAPIError_Forbidden(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
	}
	err := wrapAPIError(ghErr)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestWrapAPIError_RateLimit(t *testing.T) {
	ghErr := &gh.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Request:    &http.Request{URL: &url.URL{}},
		},
	}
	err := wrapAPIError(ghErr)
	assert.Contains(t, err.Error(), "rate limit exceeded")
}

func TestWrapAPIError_GenericError(t *testing.T) {
	err := wrapAPIError(fmt.Errorf("network timeout"))
	assert.Contains(t, err.Error(), "GitHub API error")
	assert.Contains(t, err.Error(), "network timeout")
}
