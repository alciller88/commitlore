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
	// The author "<script>" must be HTML-escaped; it must never appear unescaped
	// as a raw HTML tag (bare "<script>" only from our own Chart.js blocks is acceptable).
	assert.Contains(t, out, "&lt;script&gt;")
	assert.NotContains(t, out, `class="cl-author"><script>`)
	assert.NotContains(t, out, `class="cl-author"><script>`)
}

func TestRender_htmlUsesThemeColors(t *testing.T) {
	s := loadTestStyle(t, "epic")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "#C9A84C")
	assert.Contains(t, out, "#0F0A05")
}

func TestRender_htmlUsesCustomCSS(t *testing.T) {
	s := loadTestStyle(t, "formal")
	s.Theme.CustomCSS = "body { letter-spacing: 0.02em; }"
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "letter-spacing")
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
	assert.Contains(t, out, "Repository")
}

func TestRenderStory_htmlContainsSections(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "formal")
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "v1.0.0")
	assert.Contains(t, out, "Alice")
	assert.Contains(t, out, "2025-01")
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
	assert.Contains(t, out, "narrative")
}

func TestRender_htmlIncludesNarrativeContent(t *testing.T) {
	s := loadTestStyle(t, "formal")
	narrative := "CHANGELOG REPORT\n\nNew features have been added this sprint."
	out, err := Render(narrative, sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "CHANGELOG REPORT")
	assert.Contains(t, out, "New features have been added")
	assert.Contains(t, out, "narrative")
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

func TestRenderChangelog_themeOverride(t *testing.T) {
	s := loadTestStyle(t, "formal")
	override := &HTMLTheme{
		Background: "#FF0000",
		Primary:    "#00FF00",
		Text:       "#0000FF",
	}
	out, err := RenderWithTheme("", sampleChangelog(), s, FormatHTML, override)
	require.NoError(t, err)
	assert.Contains(t, out, "#FF0000")
	assert.Contains(t, out, "#00FF00")
	assert.Contains(t, out, "#0000FF")
}

func TestRenderStory_themeOverride(t *testing.T) {
	ch := sampleStoryChronology()
	s := loadTestStyle(t, "formal")
	override := &HTMLTheme{
		Background: "#AA0000",
		Primary:    "#00AA00",
	}
	out, err := RenderStoryWithTheme("", ch, s, FormatHTML, override)
	require.NoError(t, err)
	assert.Contains(t, out, "#AA0000")
	assert.Contains(t, out, "#00AA00")
}

func TestRenderWithTheme_nilOverride(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out1, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	out2, err := RenderWithTheme("", sampleChangelog(), s, FormatHTML, nil)
	require.NoError(t, err)
	assert.Equal(t, out1, out2)
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

func TestRenderChangelog_usesHTMLTemplateChangelog(t *testing.T) {
	s := loadTestStyle(t, "formal")
	assert.NotEmpty(t, s.HTMLTemplateChangelog, "formal style must have html_template_changelog set")
	out, err := Render("Test narrative", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "Chart.js")
	assert.Contains(t, out, "add login")
	assert.Contains(t, out, "Alice")
}

func TestRenderChangelog_fallsBackToDefault(t *testing.T) {
	s := loadTestStyle(t, "formal")
	s.HTMLTemplateChangelog = ""
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "Changelog")
	assert.Contains(t, out, "Features")
	assert.NotContains(t, out, "Chart.js")
}

func TestRenderStory_usesHTMLTemplateStory(t *testing.T) {
	s := loadTestStyle(t, "formal")
	assert.NotEmpty(t, s.HTMLTemplateStory, "formal style must have html_template_story set")
	ch := sampleStoryChronology()
	out, err := RenderStory("Story narrative", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "Chart.js")
	assert.Contains(t, out, "Alice")
}

func TestRenderStory_fallsBackToDefault(t *testing.T) {
	s := loadTestStyle(t, "formal")
	s.HTMLTemplateStory = ""
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
	assert.Contains(t, out, "Repository Story")
	assert.NotContains(t, out, "Chart.js")
}

func TestHTMLTemplateContext_allFieldsPopulated(t *testing.T) {
	s := loadTestStyle(t, "formal")
	cl := sampleChangelog()
	ctx := buildChangelogContext("some narrative", cl, s)
	assert.Equal(t, "Changelog", ctx.Title)
	assert.NotEmpty(t, ctx.Content)
	assert.NotEmpty(t, ctx.Items)
	assert.NotEmpty(t, ctx.Theme.Colors.Primary)
	assert.NotEmpty(t, ctx.Theme.Colors.Background)
	assert.NotEmpty(t, ctx.Theme.Typography.FontFamily)
	assert.NotEmpty(t, ctx.Generated)

	ch := sampleStoryChronology()
	sCtx := buildStoryContext("story text", ch, s)
	assert.Equal(t, "Repository Story", sCtx.Title)
	assert.NotEmpty(t, sCtx.Content)
	assert.Equal(t, 62, sCtx.TotalCommits)
	assert.Equal(t, "Alice", sCtx.FirstAuthor)
	assert.NotEmpty(t, sCtx.FirstDate)
	assert.Len(t, sCtx.Tags, 1)
	assert.Len(t, sCtx.Peaks, 2)
	assert.Len(t, sCtx.Contributors, 2)
}

func TestHTMLTemplateContext_repoNamePopulated(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ch := sampleStoryChronology()
	ctx := buildStoryContext("narrative", ch, s)
	assert.NotEmpty(t, ctx.RepoName)
	assert.Equal(t, "Repository", ctx.RepoName)
}

func TestHTMLTemplateContext_itemTypesInOutput(t *testing.T) {
	s := loadTestStyle(t, "formal")
	cl := sampleChangelog()
	out, err := Render("", cl, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, `type: "feat"`)
	assert.Contains(t, out, `type: "fix"`)
}

func TestHTMLTemplateContext_changelogItemsPopulated(t *testing.T) {
	s := loadTestStyle(t, "formal")
	cl := sampleChangelog()
	ctx := buildChangelogContext("", cl, s)
	assert.Len(t, ctx.Items, 3)
	typeMap := make(map[string]int)
	for _, item := range ctx.Items {
		typeMap[item.Type]++
	}
	assert.Equal(t, 2, typeMap["feat"])
	assert.Equal(t, 1, typeMap["fix"])
}
