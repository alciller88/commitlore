package main

import (
	"sync"

	"github.com/alciller88/commitlore/internal/git"
)

type commitCache struct {
	mu      sync.Mutex
	ref     string
	commits []git.Commit
}

var globalCommitCache = &commitCache{}

func (c *commitCache) get(ref string) ([]git.Commit, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ref == ref && c.commits != nil {
		return c.commits, true
	}
	return nil, false
}

func (c *commitCache) set(ref string, commits []git.Commit) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ref = ref
	c.commits = commits
}

// shouldUseCache returns true when opts have no filters that would
// make cached results invalid.
func shouldUseCache(opts git.LogOptions) bool {
	return opts.Author == "" && opts.Since.IsZero() && opts.Until.IsZero()
}

// InvalidateCache clears the commit cache. Call when the user selects
// a different repo in the Dashboard.
func InvalidateCache() {
	globalCommitCache.set("", nil)
}
