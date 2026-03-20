package git

import (
	"fmt"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CommitFiles returns the list of file paths changed in a specific commit.
func (r *Repo) CommitFiles(hash string) ([]string, error) {
	commit, err := resolveCommit(r.repo, hash)
	if err != nil {
		return nil, err
	}
	return changedFiles(commit)
}

func resolveCommit(repo *gogit.Repository, hash string) (*object.Commit, error) {
	h, err := repo.ResolveRevision(plumbing.Revision(hash))
	if err != nil {
		return nil, fmt.Errorf("resolving commit %s: %w", hash, err)
	}
	commit, err := repo.CommitObject(*h)
	if err != nil {
		return nil, fmt.Errorf("reading commit %s: %w", hash, err)
	}
	return commit, nil
}

func changedFiles(commit *object.Commit) ([]string, error) {
	parent, err := parentTree(commit)
	if err != nil {
		return nil, err
	}
	currentTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("reading commit tree: %w", err)
	}
	return diffTreeFiles(parent, currentTree)
}

func parentTree(commit *object.Commit) (*object.Tree, error) {
	if commit.NumParents() == 0 {
		return nil, nil
	}
	parent, err := commit.Parents().Next()
	if err != nil {
		return nil, fmt.Errorf("reading parent commit: %w", err)
	}
	tree, err := parent.Tree()
	if err != nil {
		return nil, fmt.Errorf("reading parent tree: %w", err)
	}
	return tree, nil
}

func diffTreeFiles(from, to *object.Tree) ([]string, error) {
	changes, err := object.DiffTree(from, to)
	if err != nil {
		return nil, fmt.Errorf("computing diff: %w", err)
	}
	files := make([]string, 0, len(changes))
	for _, c := range changes {
		name := changeName(c)
		if name != "" {
			files = append(files, name)
		}
	}
	return files, nil
}

func changeName(c *object.Change) string {
	if c.To.Name != "" {
		return c.To.Name
	}
	return c.From.Name
}
