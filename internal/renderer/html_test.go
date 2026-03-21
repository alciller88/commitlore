package renderer

import (
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_htmlChangelog(t *testing.T) {
	out, err := Render("", sampleChangelog(), FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "<html")
	assert.Contains(t, out, "<body>")
	assert.Contains(t, out, "</html>")
}

func TestRender_htmlContainsContent(t *testing.T) {
	out, err := Render("", sampleChangelog(), FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "Changelog")
	assert.Contains(t, out, "Features")
	assert.Contains(t, out, "add login")
	assert.Contains(t, out, "Alice")
}

func TestRender_htmlCommitBadges(t *testing.T) {
	out, err := Render("", sampleChangelog(), FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "type-feat")
	assert.Contains(t, out, "type-fix")
}

func TestRender_htmlEscapesContent(t *testing.T) {
	cl := changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{
				Type: changelog.TypeFeat,
				Commits: []changelog.ParsedCommit{
					{Hash: "abc1234567", Author: "<script>", Message: "feat: <b>xss</b>",
						Date: time.Now(), Type: changelog.TypeFeat},
				},
			},
		},
	}
	out, err := Render("", cl, FormatHTML)
	require.NoError(t, err)
	assert.NotContains(t, out, "<script>")
	assert.Contains(t, out, "&lt;script&gt;")
}

func TestRenderStory_html(t *testing.T) {
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "<body>")
	assert.Contains(t, out, "Repository Story")
}

func TestRenderStory_htmlContainsSections(t *testing.T) {
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "Milestones")
	assert.Contains(t, out, "Activity Peaks")
	assert.Contains(t, out, "Contributors")
	assert.Contains(t, out, "v1.0.0")
	assert.Contains(t, out, "Alice")
}

func TestRenderStory_htmlEmptyChronology(t *testing.T) {
	out, err := RenderStory("", git.Chronology{}, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<html")
	assert.NotContains(t, out, "Milestones")
}

func sampleStoryChronology() git.Chronology {
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	return git.Chronology{
		FirstCommit: git.Commit{
			Hash: "abc1234567", Author: "Alice", Email: "a@x.com",
			Date: base, Message: "initial commit",
		},
		Tags: []git.Tag{
			{Name: "v1.0.0", Hash: "def7890123", Date: base.AddDate(0, 1, 0)},
		},
		Peaks: []git.ActivityPeak{
			{Month: "2025-01", Count: 42},
			{Month: "2025-02", Count: 20},
		},
		Contributors: []git.ContributorEntry{
			{Name: "Alice", Email: "a@x.com", Date: base},
			{Name: "Bob", Email: "b@x.com", Date: base.AddDate(0, 0, 7)},
		},
		TotalCommits: 62,
	}
}
