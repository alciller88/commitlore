package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCommits_BasicFetch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{
				"sha": "abc123def456",
				"commit": {
					"message": "feat: add login\n\nDetailed description",
					"author": {
						"name": "Alice",
						"email": "alice@example.com",
						"date": "2026-01-15T10:00:00Z"
					}
				}
			},
			{
				"sha": "def789abc012",
				"commit": {
					"message": "fix: resolve crash",
					"author": {
						"name": "Bob",
						"email": "bob@example.com",
						"date": "2026-01-14T09:00:00Z"
					}
				}
			}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	commits, err := client.ListCommits(context.Background(), git.LogOptions{})
	require.NoError(t, err)
	assert.Len(t, commits, 2)

	assert.Equal(t, "abc123def456", commits[0].Hash)
	assert.Equal(t, "Alice", commits[0].Author)
	assert.Equal(t, "alice@example.com", commits[0].Email)
	assert.Equal(t, "feat: add login", commits[0].Message)

	assert.Equal(t, "def789abc012", commits[1].Hash)
	assert.Equal(t, "Bob", commits[1].Author)
}

func TestListCommits_WithLimit(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{"sha": "aaa", "commit": {"message": "first", "author": {"name": "A", "email": "a@x.com", "date": "2026-01-03T00:00:00Z"}}},
			{"sha": "bbb", "commit": {"message": "second", "author": {"name": "B", "email": "b@x.com", "date": "2026-01-02T00:00:00Z"}}},
			{"sha": "ccc", "commit": {"message": "third", "author": {"name": "C", "email": "c@x.com", "date": "2026-01-01T00:00:00Z"}}}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	commits, err := client.ListCommits(context.Background(), git.LogOptions{Limit: 2})
	require.NoError(t, err)
	assert.Len(t, commits, 2)
}

func TestListCommits_WithAuthorFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{"sha": "aaa", "commit": {"message": "by alice", "author": {"name": "Alice", "email": "alice@x.com", "date": "2026-01-02T00:00:00Z"}}},
			{"sha": "bbb", "commit": {"message": "by bob", "author": {"name": "Bob", "email": "bob@x.com", "date": "2026-01-01T00:00:00Z"}}}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	commits, err := client.ListCommits(context.Background(), git.LogOptions{Author: "alice"})
	require.NoError(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "Alice", commits[0].Author)
}

func TestListCommits_WithDateFilters(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		since := r.URL.Query().Get("since")
		until := r.URL.Query().Get("until")
		assert.NotEmpty(t, since)
		assert.NotEmpty(t, until)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[
			{"sha": "aaa", "commit": {"message": "in range", "author": {"name": "A", "email": "a@x.com", "date": "2026-01-15T00:00:00Z"}}}
		]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	since, _ := time.Parse("2006-01-02", "2026-01-01")
	until, _ := time.Parse("2006-01-02", "2026-01-31")

	commits, err := client.ListCommits(context.Background(), git.LogOptions{
		Since: since,
		Until: until,
	})
	require.NoError(t, err)
	assert.Len(t, commits, 1)
}

func TestListCommits_EmptyResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	commits, err := client.ListCommits(context.Background(), git.LogOptions{})
	require.NoError(t, err)
	assert.Empty(t, commits)
}

func TestListCommits_APIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, `{"message": "Not Found"}`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	_, err := client.ListCommits(context.Background(), git.LogOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository not found")
}

func TestListCommits_SubjectLineExtraction(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/testowner/testrepo/commits", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `[{
			"sha": "aaa",
			"commit": {
				"message": "feat: subject line\n\nThis is the body\nwith multiple lines",
				"author": {"name": "A", "email": "a@x.com", "date": "2026-01-01T00:00:00Z"}
			}
		}]`)
	})

	client, server := newTestClient(t, mux)
	defer server.Close()

	commits, err := client.ListCommits(context.Background(), git.LogOptions{})
	require.NoError(t, err)
	assert.Equal(t, "feat: subject line", commits[0].Message)
}

func TestMatchesAuthorFilter(t *testing.T) {
	c := git.Commit{Author: "Alice Smith", Email: "alice@example.com"}

	assert.True(t, matchesAuthorFilter(c, ""))
	assert.True(t, matchesAuthorFilter(c, "alice"))
	assert.True(t, matchesAuthorFilter(c, "Alice"))
	assert.True(t, matchesAuthorFilter(c, "example.com"))
	assert.False(t, matchesAuthorFilter(c, "bob"))
}

func TestSubjectLine(t *testing.T) {
	assert.Equal(t, "first line", subjectLine("first line\nsecond line"))
	assert.Equal(t, "only line", subjectLine("only line"))
	assert.Equal(t, "trimmed", subjectLine("  trimmed  "))
	assert.Equal(t, "", subjectLine(""))
}
