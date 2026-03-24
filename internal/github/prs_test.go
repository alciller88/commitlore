package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPullRequests_FetchesMergedPRs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "closed", r.URL.Query().Get("state"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{
				"number": 42,
				"title": "Add feature X",
				"state": "closed",
				"user": {"login": "alice"},
				"created_at": "2026-01-10T00:00:00Z",
				"merged_at": "2026-01-15T00:00:00Z",
				"labels": [{"name": "enhancement"}]
			},
			{
				"number": 43,
				"title": "Not merged PR",
				"state": "closed",
				"user": {"login": "bob"},
				"created_at": "2026-01-12T00:00:00Z",
				"labels": []
			}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	prs, err := client.ListPullRequests(context.Background(), time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Len(t, prs, 1)

	assert.Equal(t, 42, prs[0].Number)
	assert.Equal(t, "Add feature X", prs[0].Title)
	assert.Equal(t, "alice", prs[0].Author)
	assert.Equal(t, "closed", prs[0].State)
	assert.Equal(t, []string{"enhancement"}, prs[0].Labels)
}

func TestListPullRequests_WithSinceFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{
				"number": 10,
				"title": "Recent PR",
				"state": "closed",
				"user": {"login": "alice"},
				"created_at": "2026-02-01T00:00:00Z",
				"merged_at": "2026-02-05T00:00:00Z",
				"labels": []
			},
			{
				"number": 5,
				"title": "Old PR",
				"state": "closed",
				"user": {"login": "bob"},
				"created_at": "2025-06-01T00:00:00Z",
				"merged_at": "2025-06-10T00:00:00Z",
				"labels": []
			}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	since, _ := time.Parse("2006-01-02", "2026-01-01")
	prs, err := client.ListPullRequests(context.Background(), since, time.Time{})
	require.NoError(t, err)
	assert.Len(t, prs, 1)
	assert.Equal(t, 10, prs[0].Number)
}

func TestListPullRequests_WithUntilFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{
				"number": 20,
				"title": "Future PR",
				"state": "closed",
				"user": {"login": "alice"},
				"created_at": "2026-06-01T00:00:00Z",
				"merged_at": "2026-06-15T00:00:00Z",
				"labels": []
			},
			{
				"number": 15,
				"title": "In range PR",
				"state": "closed",
				"user": {"login": "bob"},
				"created_at": "2026-01-01T00:00:00Z",
				"merged_at": "2026-01-20T00:00:00Z",
				"labels": []
			}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	until, _ := time.Parse("2006-01-02", "2026-03-01")
	prs, err := client.ListPullRequests(context.Background(), time.Time{}, until)
	require.NoError(t, err)
	assert.Len(t, prs, 1)
	assert.Equal(t, 15, prs[0].Number)
}

func TestListPullRequests_EmptyResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	prs, err := client.ListPullRequests(context.Background(), time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Empty(t, prs)
}

func TestListPullRequests_APIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, `{"message": "API rate limit exceeded"}`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	_, err := client.ListPullRequests(context.Background(), time.Time{}, time.Time{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
}

func TestListPullRequests_OnlyGETMethod(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/pulls", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "only GET requests are allowed")
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	_, err := client.ListPullRequests(context.Background(), time.Time{}, time.Time{})
	require.NoError(t, err)
}
