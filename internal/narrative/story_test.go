package narrative

import (
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleChronology() git.Chronology {
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	return git.Chronology{
		FirstCommit: git.Commit{
			Hash:    "abc1234567",
			Author:  "Alice",
			Email:   "alice@x.com",
			Date:    base,
			Message: "initial commit",
		},
		Tags: []git.Tag{
			{Name: "v1.0.0", Hash: "def7890", Date: base.AddDate(0, 1, 0)},
			{Name: "v2.0.0", Hash: "ghi4567", Date: base.AddDate(0, 3, 0)},
		},
		Peaks: []git.ActivityPeak{
			{Month: "2025-03", Count: 42},
			{Month: "2025-01", Count: 30},
			{Month: "2025-02", Count: 20},
		},
		Contributors: []git.ContributorEntry{
			{Name: "Alice", Email: "alice@x.com", Date: base},
			{Name: "Bob", Email: "bob@x.com", Date: base.AddDate(0, 0, 7)},
		},
		TotalCommits: 92,
	}
}

func TestGenerateStory_formal(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "Repository Report")
	assert.Contains(t, out, "2025-01-01")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "92")
}

func TestGenerateStory_patchnotes(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "patchnotes"))
	require.NoError(t, err)
	assert.Contains(t, out, "LEGEND BEGINS")
	assert.Contains(t, out, "ACHIEVEMENT UNLOCKED")
}

func TestGenerateStory_ironic(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "ironic"))
	require.NoError(t, err)
	assert.Contains(t, out, "nobody asked for")
	assert.Contains(t, out, "One of us now")
}

func TestGenerateStory_epic(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "epic"))
	require.NoError(t, err)
	assert.Contains(t, out, "SAGA BEGINS")
	assert.Contains(t, out, "fellowship")
}

func TestGenerateStory_includesMilestones(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "v1.0.0")
	assert.Contains(t, out, "v2.0.0")
	assert.Contains(t, out, "Milestones")
}

func TestGenerateStory_includesPeaks(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "Activity Peaks")
	assert.Contains(t, out, "2025-03")
	assert.Contains(t, out, "42")
}

func TestGenerateStory_includesContributors(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "Contributors")
	assert.Contains(t, out, "Bob")
}

func TestGenerateStory_emptyChronology(t *testing.T) {
	out, err := GenerateStory(git.Chronology{}, loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestGenerateStory_noTags(t *testing.T) {
	ch := sampleChronology()
	ch.Tags = nil
	out, err := GenerateStory(ch, loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.NotContains(t, out, "Milestones")
}

func TestGenerateStory_footer(t *testing.T) {
	out, err := GenerateStory(sampleChronology(), loadStyle(t, "epic"))
	require.NoError(t, err)
	assert.Contains(t, out, "next chapter awaits")
}
