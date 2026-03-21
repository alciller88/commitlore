package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Tags returns all tags in the repository with their commit hash and date.
func (r *Repo) Tags() ([]Tag, error) {
	refs, err := r.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("reading tags: %w", err)
	}
	defer refs.Close()

	var tags []Tag
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		tag, resolveErr := r.resolveTag(ref)
		if resolveErr != nil {
			return nil
		}
		tags = append(tags, tag)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("iterating tags: %w", err)
	}
	return tags, nil
}

func (r *Repo) resolveTag(ref *plumbing.Reference) (Tag, error) {
	name := tagName(ref)

	commit, err := r.resolveAnnotatedTag(ref)
	if err == nil {
		return Tag{Name: name, Hash: commit.Hash.String(), Date: commit.Author.When}, nil
	}

	return r.resolveLightweightTag(ref, name)
}

func (r *Repo) resolveAnnotatedTag(ref *plumbing.Reference) (*object.Commit, error) {
	tagObj, err := r.repo.TagObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	return tagObj.Commit()
}

func (r *Repo) resolveLightweightTag(ref *plumbing.Reference, name string) (Tag, error) {
	commit, err := r.repo.CommitObject(ref.Hash())
	if err != nil {
		return Tag{}, fmt.Errorf("resolving tag %q: %w", name, err)
	}
	return Tag{Name: name, Hash: commit.Hash.String(), Date: commit.Author.When}, nil
}

func tagName(ref *plumbing.Reference) string {
	return strings.TrimPrefix(ref.Name().String(), "refs/tags/")
}
