package git

import (
	"fmt"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Repo wraps a go-git repository for read operations.
type Repo struct {
	repo *gogit.Repository
}

// Open opens a git repository at the given filesystem path.
func Open(path string) (*Repo, error) {
	r, err := gogit.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("opening repo at %s: %w", path, err)
	}
	return &Repo{repo: r}, nil
}

// Log returns commits from the repository filtered by the given options.
func (r *Repo) Log(opts LogOptions) ([]Commit, error) {
	logOpts := &gogit.LogOptions{Order: gogit.LogOrderCommitterTime}

	iter, err := r.repo.Log(logOpts)
	if err != nil {
		return nil, fmt.Errorf("reading commit log: %w", err)
	}
	defer iter.Close()

	return collectCommits(iter, opts)
}

func collectCommits(iter object.CommitIter, opts LogOptions) ([]Commit, error) {
	var commits []Commit

	err := iter.ForEach(func(c *object.Commit) error {
		if opts.Limit > 0 && len(commits) >= opts.Limit {
			return fmt.Errorf("limit reached")
		}
		if !matchesFilters(c, opts) {
			return nil
		}
		commits = append(commits, toCommit(c))
		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return nil, fmt.Errorf("iterating commits: %w", err)
	}
	return commits, nil
}

func matchesFilters(c *object.Commit, opts LogOptions) bool {
	if !matchesAuthor(c, opts.Author) {
		return false
	}
	if !opts.Since.IsZero() && c.Author.When.Before(opts.Since) {
		return false
	}
	if !opts.Until.IsZero() && c.Author.When.After(opts.Until) {
		return false
	}
	return true
}

func matchesAuthor(c *object.Commit, author string) bool {
	if author == "" {
		return true
	}
	lower := strings.ToLower(author)
	return strings.Contains(strings.ToLower(c.Author.Name), lower) ||
		strings.Contains(strings.ToLower(c.Author.Email), lower)
}

func toCommit(c *object.Commit) Commit {
	return Commit{
		Hash:    c.Hash.String(),
		Author:  c.Author.Name,
		Email:   c.Author.Email,
		Date:    c.Author.When,
		Message: strings.TrimSpace(c.Message),
	}
}
