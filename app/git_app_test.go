package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommitJSON_fieldNames(t *testing.T) {
	commit := git.Commit{
		Hash:    "abc1234567890",
		Author:  "Test Author",
		Email:   "test@example.com",
		Date:    time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC),
		Message: "feat: test commit",
	}

	data, err := json.Marshal(commit)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &m))

	assert.Contains(t, m, "hash")
	assert.Contains(t, m, "author")
	assert.Contains(t, m, "email")
	assert.Contains(t, m, "date")
	assert.Contains(t, m, "message")

	assert.Equal(t, "abc1234567890", m["hash"])
	assert.Equal(t, "Test Author", m["author"])
	assert.Equal(t, "test@example.com", m["email"])
	assert.Equal(t, "feat: test commit", m["message"])
}

func TestCommitJSON_arrayParsing(t *testing.T) {
	commits := []git.Commit{
		{
			Hash:    "abc123",
			Author:  "Alice",
			Email:   "alice@test.com",
			Date:    time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC),
			Message: "first",
		},
		{
			Hash:    "def456",
			Author:  "Bob",
			Email:   "bob@test.com",
			Date:    time.Date(2026, 3, 20, 12, 0, 0, 0, time.UTC),
			Message: "second",
		},
	}

	jsonStr, err := toJSON(commits)
	require.NoError(t, err)

	var parsed []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(jsonStr), &parsed))

	assert.Len(t, parsed, 2)
	assert.Equal(t, "abc123", parsed[0]["hash"])
	assert.Equal(t, "Alice", parsed[0]["author"])
	assert.Equal(t, "def456", parsed[1]["hash"])
}
