package git

import (
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTags_empty(t *testing.T) {
	dir := createTestRepo(t, sampleCommits()[:1])
	repo, err := Open(dir)
	require.NoError(t, err)

	tags, err := repo.Tags()
	assert.NoError(t, err)
	assert.Empty(t, tags)
}

func TestTags_lightweightTag(t *testing.T) {
	dir := createTestRepo(t, sampleCommits()[:1])
	createLightweightTag(t, dir, "v0.1.0")

	repo, err := Open(dir)
	require.NoError(t, err)

	tags, err := repo.Tags()
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, "v0.1.0", tags[0].Name)
	assert.NotEmpty(t, tags[0].Hash)
	assert.False(t, tags[0].Date.IsZero())
}

func TestTags_multipleTags(t *testing.T) {
	commits := sampleCommits()
	dir := createTestRepo(t, commits)
	createLightweightTag(t, dir, "v0.1.0")
	createLightweightTag(t, dir, "v0.2.0")

	repo, err := Open(dir)
	require.NoError(t, err)

	tags, err := repo.Tags()
	assert.NoError(t, err)
	assert.Len(t, tags, 2)
}

func createLightweightTag(t *testing.T, dir, name string) {
	t.Helper()
	r, err := gogit.PlainOpen(dir)
	require.NoError(t, err)

	head, err := r.Head()
	require.NoError(t, err)

	ref := plumbing.NewHashReference(
		plumbing.NewTagReferenceName(name),
		head.Hash(),
	)
	err = r.Storer.SetReference(ref)
	require.NoError(t, err)
}

func TestTags_annotatedTag(t *testing.T) {
	commits := sampleCommits()[:1]
	dir := createTestRepo(t, commits)
	createAnnotatedTag(t, dir, "v1.0.0")

	repo, err := Open(dir)
	require.NoError(t, err)

	tags, err := repo.Tags()
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, "v1.0.0", tags[0].Name)
}

func createAnnotatedTag(t *testing.T, dir, name string) {
	t.Helper()
	r, err := gogit.PlainOpen(dir)
	require.NoError(t, err)

	head, err := r.Head()
	require.NoError(t, err)

	commit, err := r.CommitObject(head.Hash())
	require.NoError(t, err)

	_, err = r.CreateTag(name, commit.Hash, &gogit.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Tagger",
			Email: "tagger@example.com",
			When:  time.Now(),
		},
		Message: "release " + name,
	})
	require.NoError(t, err)
}
