package git

import (
	"os"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestRepo(t *testing.T, commits []testCommit) string {
	t.Helper()
	dir := t.TempDir()

	r, err := gogit.PlainInit(dir, false)
	require.NoError(t, err)

	wt, err := r.Worktree()
	require.NoError(t, err)

	for _, tc := range commits {
		writeFile(t, dir, tc)
		err = wt.AddGlob(".")
		require.NoError(t, err)
		_, err = wt.Commit(tc.message, &gogit.CommitOptions{
			Author: &object.Signature{
				Name:  tc.author,
				Email: tc.email,
				When:  tc.when,
			},
		})
		require.NoError(t, err)
	}

	return dir
}

type testCommit struct {
	author  string
	email   string
	message string
	when    time.Time
	file    string
}

func writeFile(t *testing.T, dir string, tc testCommit) {
	t.Helper()
	name := tc.file
	if name == "" {
		name = "file.txt"
	}
	path := dir + "/" + name
	err := os.WriteFile(path, []byte(tc.message), 0644)
	require.NoError(t, err)
}

func sampleCommits() []testCommit {
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	return []testCommit{
		{author: "Alice", email: "alice@example.com", message: "first commit", when: base, file: "a.txt"},
		{author: "Bob", email: "bob@example.com", message: "second commit", when: base.Add(time.Hour), file: "b.txt"},
		{author: "Alice", email: "alice@example.com", message: "third commit", when: base.Add(2 * time.Hour), file: "c.txt"},
		{author: "Charlie", email: "charlie@example.com", message: "fourth commit", when: base.Add(3 * time.Hour), file: "d.txt"},
	}
}

func TestOpen_validRepo(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
}

func TestOpen_invalidPath(t *testing.T) {
	repo, err := Open("/nonexistent/path")
	assert.Error(t, err)
	assert.Nil(t, repo)
}

func TestLog_allCommits(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	commits, err := repo.Log(LogOptions{})
	assert.NoError(t, err)
	assert.Len(t, commits, 4)
}

func TestLog_filterByAuthorName(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	commits, err := repo.Log(LogOptions{Author: "Alice"})
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
	for _, c := range commits {
		assert.Equal(t, "Alice", c.Author)
	}
}

func TestLog_filterByAuthorEmail(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	commits, err := repo.Log(LogOptions{Author: "bob@example.com"})
	assert.NoError(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "Bob", commits[0].Author)
}

func TestLog_filterByAuthorCaseInsensitive(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	commits, err := repo.Log(LogOptions{Author: "alice"})
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
}

func TestLog_filterBySince(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	since := time.Date(2025, 1, 1, 13, 30, 0, 0, time.UTC)
	commits, err := repo.Log(LogOptions{Since: since})
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
}

func TestLog_filterByUntil(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	until := time.Date(2025, 1, 1, 13, 30, 0, 0, time.UTC)
	commits, err := repo.Log(LogOptions{Until: until})
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
}

func TestLog_filterByLimit(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	commits, err := repo.Log(LogOptions{Limit: 2})
	assert.NoError(t, err)
	assert.Len(t, commits, 2)
}

func TestLog_combinedFilters(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	since := time.Date(2025, 1, 1, 12, 30, 0, 0, time.UTC)
	commits, err := repo.Log(LogOptions{
		Author: "Alice",
		Since:  since,
		Limit:  10,
	})
	assert.NoError(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, "third commit", commits[0].Message)
}

func TestLog_commitFields(t *testing.T) {
	dir := createTestRepo(t, sampleCommits())
	repo, err := Open(dir)
	require.NoError(t, err)

	commits, err := repo.Log(LogOptions{Limit: 1})
	require.NoError(t, err)
	require.Len(t, commits, 1)

	c := commits[0]
	assert.NotEmpty(t, c.Hash)
	assert.Len(t, c.Hash, 40)
	assert.NotEmpty(t, c.Author)
	assert.NotEmpty(t, c.Email)
	assert.NotEmpty(t, c.Message)
	assert.False(t, c.Date.IsZero())
}

func TestLog_emptyRepo(t *testing.T) {
	dir := t.TempDir()
	_, err := gogit.PlainInit(dir, false)
	require.NoError(t, err)

	repo, err := Open(dir)
	require.NoError(t, err)

	_, err = repo.Log(LogOptions{})
	assert.Error(t, err)
}
