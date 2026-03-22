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
	s := loadTestStyle(t, "formal")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "<html")
	assert.Contains(t, out, "<body>")
	assert.Contains(t, out, "</html>")
}

func TestRender_htmlContainsContent(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "Changelog")
	assert.Contains(t, out, "Features")
	assert.Contains(t, out, "add login")
	assert.Contains(t, out, "Alice")
}

func TestRender_htmlCommitBadges(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
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
	s := loadTestStyle(t, "formal")
	out, err := Render("", cl, s, FormatHTML)
	require.NoError(t, err)
	assert.NotContains(t, out, "<script>")
	assert.Contains(t, out, "&lt;script&gt;")
}

func TestRender_htmlUsesThemeColors(t *testing.T) {
	s := loadTestStyle(t, "epic")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "#D4AF37")
	assert.Contains(t, out, "#1A1209")
}

func TestRender_htmlUsesCustomCSS(t *testing.T) {
	s := loadTestStyle(t, "epic")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "text-shadow")
}

func TestRender_htmlBorderedCardStyle(t *testing.T) {
	s := loadTestStyle(t, "patchnotes")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "border-radius: 4px")
}

func TestRenderStory_html(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "formal")
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "<body>")
	assert.Contains(t, out, "Repository Story")
}

func TestRenderStory_htmlContainsSections(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "formal")
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "Milestones")
	assert.Contains(t, out, "Activity Peaks")
	assert.Contains(t, out, "Contributors")
	assert.Contains(t, out, "v1.0.0")
	assert.Contains(t, out, "Alice")
}

func TestRenderStory_htmlEmptyChronology(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out, err := RenderStory("", git.Chronology{}, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<html")
	assert.NotContains(t, out, "Milestones")
}

func TestRenderStory_htmlIncludesNarrativeContent(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "epic")
	narrative := "THE SAGA BEGINS\n\nIn the age of legends, the great Alice laid the first stone."
	out, err := RenderStory(narrative, ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "THE SAGA BEGINS")
	assert.Contains(t, out, "Alice laid the first stone")
	assert.Contains(t, out, "class=\"narrative\"")
}

func TestRender_htmlIncludesNarrativeContent(t *testing.T) {
	s := loadTestStyle(t, "formal")
	narrative := "CHANGELOG REPORT\n\nNew features have been added this sprint."
	out, err := Render(narrative, sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "CHANGELOG REPORT")
	assert.Contains(t, out, "New features have been added")
	assert.Contains(t, out, "class=\"narrative\"")
}

func TestRenderStory_htmlRendersMarkdown(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "formal")
	narrative := "## The Beginning\n\nThe project started with a **bold** vision.\n\n* First item\n* Second item"
	out, err := RenderStory(narrative, ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<h2>The Beginning</h2>")
	assert.Contains(t, out, "<strong>bold</strong>")
	assert.Contains(t, out, "<li>First item</li>")
	assert.NotContains(t, out, "## The Beginning")
}

func TestRenderChangelog_htmlRendersMarkdown(t *testing.T) {
	s := loadTestStyle(t, "formal")
	narrative := "## Added\n\n* New login feature\n* Dashboard redesign"
	out, err := Render(narrative, sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<h2>Added</h2>")
	assert.Contains(t, out, "<li>New login feature</li>")
	assert.NotContains(t, out, "## Added")
}

func TestRenderStory_htmlNarrativeEscapesHTML(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "formal")
	narrative := "<script>alert('xss')</script>"
	out, err := RenderStory(narrative, ch, s, FormatHTML)
	require.NoError(t, err)
	assert.NotContains(t, out, "<script>alert")
	assert.NotContains(t, out, "alert('xss')")
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
