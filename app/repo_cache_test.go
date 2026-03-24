package main

import (
	"testing"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
)

func TestCommitCache_missOnEmpty(t *testing.T) {
	c := &commitCache{}
	_, ok := c.get("owner/repo")
	assert.False(t, ok)
}

func TestCommitCache_hitAfterSet(t *testing.T) {
	c := &commitCache{}
	commits := []git.Commit{{Hash: "abc", Author: "Alice"}}
	c.set("owner/repo", commits)
	got, ok := c.get("owner/repo")
	assert.True(t, ok)
	assert.Equal(t, commits, got)
}

func TestCommitCache_missOnDifferentRef(t *testing.T) {
	c := &commitCache{}
	c.set("owner/repo", []git.Commit{{Hash: "abc"}})
	_, ok := c.get("owner/repo2")
	assert.False(t, ok)
}

func TestCommitCache_invalidate(t *testing.T) {
	c := &commitCache{}
	c.set("owner/repo", []git.Commit{{Hash: "abc"}})
	c.set("", nil)
	_, ok := c.get("owner/repo")
	assert.False(t, ok)
}

func TestShouldUseCache_defaultOpts(t *testing.T) {
	assert.True(t, shouldUseCache(git.LogOptions{}))
}

func TestShouldUseCache_withAuthor(t *testing.T) {
	assert.False(t, shouldUseCache(git.LogOptions{Author: "alice"}))
}
