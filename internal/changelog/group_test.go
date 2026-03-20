package changelog

import (
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleCommits() []git.Commit {
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	return []git.Commit{
		{Hash: "aaa", Author: "Alice", Email: "a@x.com", Date: base, Message: "feat: add login"},
		{Hash: "bbb", Author: "Bob", Email: "b@x.com", Date: base, Message: "fix: resolve crash"},
		{Hash: "ccc", Author: "Alice", Email: "a@x.com", Date: base, Message: "feat: add logout"},
		{Hash: "ddd", Author: "Charlie", Email: "c@x.com", Date: base, Message: "chore: update deps"},
		{Hash: "eee", Author: "Alice", Email: "a@x.com", Date: base, Message: "docs: update README"},
		{Hash: "fff", Author: "Bob", Email: "b@x.com", Date: base, Message: "initial commit"},
	}
}

func TestGroupCommits_groupsByType(t *testing.T) {
	cl := GroupCommits(sampleCommits())

	typeMap := make(map[CommitType]int)
	for _, g := range cl.Groups {
		typeMap[g.Type] = len(g.Commits)
	}

	assert.Equal(t, 2, typeMap[TypeFeat])
	assert.Equal(t, 1, typeMap[TypeFix])
	assert.Equal(t, 1, typeMap[TypeChore])
	assert.Equal(t, 1, typeMap[TypeDocs])
	assert.Equal(t, 1, typeMap[TypeOther])
}

func TestGroupCommits_orderIsCorrect(t *testing.T) {
	cl := GroupCommits(sampleCommits())

	require.True(t, len(cl.Groups) >= 2)

	// feat should come before fix in the ordering
	var featIdx, fixIdx int
	for i, g := range cl.Groups {
		if g.Type == TypeFeat {
			featIdx = i
		}
		if g.Type == TypeFix {
			fixIdx = i
		}
	}
	assert.Less(t, featIdx, fixIdx)
}

func TestGroupCommits_emptyInput(t *testing.T) {
	cl := GroupCommits(nil)
	assert.Empty(t, cl.Groups)
}

func TestGroupCommits_preservesCommitData(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc123", Author: "Dev", Email: "dev@x.com", Date: time.Now(), Message: "feat: something"},
	}

	cl := GroupCommits(commits)
	require.Len(t, cl.Groups, 1)
	require.Len(t, cl.Groups[0].Commits, 1)

	pc := cl.Groups[0].Commits[0]
	assert.Equal(t, "abc123", pc.Hash)
	assert.Equal(t, "Dev", pc.Author)
	assert.Equal(t, "dev@x.com", pc.Email)
	assert.Equal(t, TypeFeat, pc.Type)
}

func TestGroupCommits_singleType(t *testing.T) {
	commits := []git.Commit{
		{Hash: "a", Message: "fix: one"},
		{Hash: "b", Message: "fix: two"},
		{Hash: "c", Message: "fix: three"},
	}

	cl := GroupCommits(commits)
	require.Len(t, cl.Groups, 1)
	assert.Equal(t, TypeFix, cl.Groups[0].Type)
	assert.Len(t, cl.Groups[0].Commits, 3)
}
