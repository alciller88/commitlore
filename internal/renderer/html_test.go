package renderer

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/styles"
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
			{Name: "Alice", Email: "a@x.com", Date: base, Count: 42},
			{Name: "Bob", Email: "b@x.com", Date: base.AddDate(0, 0, 7), Count: 20},
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
	ctx := buildChangelogContext("some narrative", cl, s, "TestRepo", "")
	assert.Equal(t, "Changelog", ctx.Title)
	assert.NotEmpty(t, ctx.Content)
	assert.NotEmpty(t, ctx.Items)
	assert.NotEmpty(t, ctx.Theme.Colors.Primary)
	assert.NotEmpty(t, ctx.Theme.Colors.Background)
	assert.NotEmpty(t, ctx.Theme.Typography.FontFamily)
	assert.NotEmpty(t, ctx.Generated)

	ch := sampleStoryChronology()
	sCtx := buildStoryContext("story text", ch, s, "TestRepo")
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
	ctx := buildStoryContext("narrative", ch, s, "TestRepo")
	assert.NotEmpty(t, ctx.RepoName)
	assert.Equal(t, "TestRepo", ctx.RepoName)
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
	ctx := buildChangelogContext("", cl, s, "TestRepo", "")
	assert.Len(t, ctx.Items, 3)
	typeMap := make(map[string]int)
	for _, item := range ctx.Items {
		typeMap[item.Type]++
	}
	assert.Equal(t, 2, typeMap["feat"])
	assert.Equal(t, 1, typeMap["fix"])
}

func TestRenderChangelog_repoNameInOutput(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out, err := Render("", sampleChangelog(), s, FormatHTML, "commitlore")
	require.NoError(t, err)
	assert.Contains(t, out, "commitlore")
}

func TestHTMLTemplateContext_commitsByWeekPopulated(t *testing.T) {
	s := loadTestStyle(t, "formal")
	cl := sampleChangelog()
	ctx := buildChangelogContext("", cl, s, "TestRepo", "")
	assert.NotEmpty(t, ctx.CommitsByWeek)
	total := 0
	for _, w := range ctx.CommitsByWeek {
		assert.NotEmpty(t, w.Label)
		assert.Greater(t, w.Count, 0)
		total += w.Count
	}
	assert.Equal(t, len(ctx.Items), total)
}

// Bug 1 — SVG logo not escaped in default renderer
func TestWriteSiteHeaderLogo_inlineSVG(t *testing.T) {
	var buf bytes.Buffer
	theme := styles.Theme{Logo: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><circle r="10"/></svg>`}
	writeSiteHeaderLogo(&buf, theme)
	assert.Contains(t, buf.String(), "<svg")
	assert.NotContains(t, buf.String(), "&lt;svg")
	assert.Contains(t, buf.String(), `class="logo-svg"`)
}

func TestWriteSiteHeaderLogo_urlLogo(t *testing.T) {
	var buf bytes.Buffer
	theme := styles.Theme{Logo: "https://example.com/logo.png"}
	writeSiteHeaderLogo(&buf, theme)
	assert.Contains(t, buf.String(), `<img class="logo" src="https://example.com/logo.png"`)
	assert.NotContains(t, buf.String(), `class="logo-svg"`)
}

// Bug 1b — safeHTML in custom templates
func TestCustomTemplate_svgLogoRendersUnescaped(t *testing.T) {
	svgLogo := `<svg xmlns="http://www.w3.org/2000/svg"><circle r="5"/></svg>`
	tmpl := `<div>{{safeHTML .Theme.Logo}}</div>`
	ctx := HTMLTemplateContext{Theme: styles.Theme{Logo: svgLogo}}
	out, err := executeHTMLTemplate(tmpl, ctx)
	require.NoError(t, err)
	assert.Contains(t, out, "<svg")
	assert.NotContains(t, out, "&lt;svg")
}

// Bug 2 — Generated is a real date
func TestBuildChangelogContext_generatedIsDate(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ctx := buildChangelogContext("", sampleChangelog(), s, "repo", "")
	assert.NotEqual(t, "Generated by CommitLore", ctx.Generated)
	assert.Regexp(t, `\d{1,2} \w+ \d{4}`, ctx.Generated)
}

func TestBuildStoryContext_generatedIsDate(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ctx := buildStoryContext("", sampleStoryChronology(), s, "repo")
	assert.NotEqual(t, "Generated by CommitLore", ctx.Generated)
	assert.Regexp(t, `\d{1,2} \w+ \d{4}`, ctx.Generated)
}

// Bug 3 — FontSizeH has a default
func TestWithDefaults_fontSizeHNotEmpty(t *testing.T) {
	theme := styles.Theme{}
	d := withDefaults(theme)
	assert.NotEmpty(t, d.Typography.FontSizeH, "FontSizeH must have a default fallback")
}

// Bug 4 — mul and divf are registered
func TestExecuteHTMLTemplate_mulDivfFunctions(t *testing.T) {
	tmpl := `{{$r := mul (divf 3 10) 100}}{{printf "%.0f" $r}}`
	out, err := executeHTMLTemplate(tmpl, struct{}{})
	require.NoError(t, err)
	assert.Equal(t, "30", out)
}

func TestExecuteHTMLTemplate_divfByZero(t *testing.T) {
	tmpl := `{{divf 5 0}}`
	out, err := executeHTMLTemplate(tmpl, struct{}{})
	require.NoError(t, err)
	assert.Equal(t, "0", out)
}

// Bug 4b — patchnotes story renders without error
func TestRenderStory_patchnotesNoError(t *testing.T) {
	s := loadTestStyle(t, "patchnotes")
	ch := sampleStoryChronology()
	out, err := RenderStory("narrative", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<!DOCTYPE html>")
}

// Bug 1 — semverFromString
func TestSemverFromString(t *testing.T) {
	cases := []struct{ input, want string }{
		{"v1.2.0", "v1.2.0"},
		{"v0.0.1", "v0.0.1"},
		{"v2.0.0-alpha", "v2.0.0-alpha"},
		{"HEAD", ""},
		{"main", ""},
		{"2025-01-01", ""},
		{"", ""},
		{"v1.2", ""},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			assert.Equal(t, c.want, SemverFromString(c.input))
		})
	}
}

// Bug 1b — version propagates through context
func TestBuildChangelogContext_versionFromSemver(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ctx := buildChangelogContext("", sampleChangelog(), s, "repo", "v1.5.0")
	assert.Equal(t, "v1.5.0", ctx.Version)
}

func TestBuildChangelogContext_versionEmptyWhenNotSemver(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ctx := buildChangelogContext("", sampleChangelog(), s, "repo", "")
	assert.Empty(t, ctx.Version)
}

// Bug 1c — version badge visible in rendered output
func TestRenderChangelog_versionBadgeVisible(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ctx := buildChangelogContext("", sampleChangelog(), s, "repo", "v2.1.0")
	out, err := executeHTMLTemplate(s.HTMLTemplateChangelog, ctx)
	require.NoError(t, err)
	assert.Contains(t, out, "v2.1.0")
}

// Bug 2 — contributor Count propagates
func TestBuildStoryContext_contributorCountPopulated(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ch := sampleStoryChronology()
	ctx := buildStoryContext("", ch, s, "repo")
	for _, c := range ctx.Contributors {
		assert.Greater(t, c.Count, 0,
			"contributor %q should have Count > 0", c.Name)
	}
}

func TestRenderStory_contributorCountInOutput(t *testing.T) {
	s := loadTestStyle(t, "patchnotes")
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	// Contributor cards must not show ">0 commits<" (zero count)
	assert.NotContains(t, out, ">0 commits<",
		"contributor count should not be zero")
	assert.Contains(t, out, ">42 commits<")
	assert.Contains(t, out, ">20 commits<")
}

// Bug 3 — FirstAuthor fallback
func TestBuildStoryContext_firstAuthorFallback(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ch := sampleStoryChronology()
	ch.FirstCommit.Author = ""
	ctx := buildStoryContext("", ch, s, "repo")
	assert.NotEmpty(t, ctx.FirstAuthor)
	assert.Equal(t, "an unknown contributor", ctx.FirstAuthor)
}

// Bug 4 — type badge CSS covers all commit types (default renderer)
func TestRender_htmlTypeBadgesAllTypes(t *testing.T) {
	s := loadTestStyle(t, "formal")
	s.HTMLTemplateChangelog = ""
	cl := changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{Type: changelog.TypeRefactor, Commits: []changelog.ParsedCommit{
				{Hash: "aaa0000001", Author: "A", Date: time.Now(),
					Type: changelog.TypeRefactor, Message: "refactor: cleanup"},
			}},
			{Type: changelog.TypeDocs, Commits: []changelog.ParsedCommit{
				{Hash: "bbb0000002", Author: "B", Date: time.Now(),
					Type: changelog.TypeDocs, Message: "docs: update readme"},
			}},
			{Type: changelog.TypeChore, Commits: []changelog.ParsedCommit{
				{Hash: "ccc0000003", Author: "C", Date: time.Now(),
					Type: changelog.TypeChore, Message: "chore: bump deps"},
			}},
		},
	}
	out, err := Render("", cl, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "type-refactor")
	assert.Contains(t, out, "type-docs")
	assert.Contains(t, out, "type-chore")
}

// Bug 4b — formal custom template type badges
func TestRender_formalHTMLTypeBadgesAllTypes(t *testing.T) {
	s := loadTestStyle(t, "formal")
	require.NotEmpty(t, s.HTMLTemplateChangelog)
	cl := changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{Type: changelog.TypeRefactor, Commits: []changelog.ParsedCommit{
				{Hash: "aaa0000001", Author: "A", Date: time.Now(),
					Type: changelog.TypeRefactor, Message: "refactor: cleanup"},
			}},
		},
	}
	out, err := Render("", cl, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "type-refactor")
}

// Phase 3 — Bug 1: formal story Generated not duplicated in banner
func TestRenderStory_formalGeneratedNotDuplicated(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	generated := time.Now().Format("2 Jan 2006")
	// Count occurrences: Generated should appear at most twice
	// (once in date range "FirstDate → Generated", once perhaps in meta)
	// but NOT as "Generated · Generated" (the old doubled pattern)
	assert.NotContains(t, out, generated+" · "+generated,
		"Generated date should not appear twice separated by middot")
}

// Phase 3 — Bug 2: patchnotes fade-in fallback timeout
func TestRenderChangelog_patchnotesHasFadeinFallback(t *testing.T) {
	s := loadTestStyle(t, "patchnotes")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "setTimeout",
		"patchnotes changelog must have setTimeout fallback for fade-in animations")
}

func TestRenderStory_patchnotesHasFadeinFallback(t *testing.T) {
	s := loadTestStyle(t, "patchnotes")
	out, err := RenderStory("", sampleStoryChronology(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "setTimeout",
		"patchnotes story must have setTimeout fallback for fade-in animations")
}

// Phase 3 — Bug 3: epic unique canvas IDs
func TestRenderChangelog_epicUniqueCanvasIDs(t *testing.T) {
	s := loadTestStyle(t, "epic")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.NotContains(t, out, `id="activityChart"`,
		"epic changelog must not use generic activityChart ID")
	assert.Contains(t, out, "epic-cl-",
		"epic changelog must use prefixed canvas IDs")
}

// Phase 3 — Bug 4: ironic title shows ironic default
func TestRenderChangelog_ironicTitleIsIronic(t *testing.T) {
	s := loadTestStyle(t, "ironic")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "changelog")
	assert.NotContains(t, out, `<div class="doc-title">Changelog</div>`,
		"ironic title should not show generic 'Changelog'")
}

func TestRenderStory_ironicTitleIsIronic(t *testing.T) {
	s := loadTestStyle(t, "ironic")
	out, err := RenderStory("", sampleStoryChronology(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "a story i guess",
		"ironic story title should be ironic default")
	assert.NotContains(t, out, `<div class="doc-title">Repository Story</div>`)
}

// Phase 3 — Bug 5: ironic story contributor table no duplicate name
func TestRenderStory_ironicContributorTableNoDuplicateName(t *testing.T) {
	s := loadTestStyle(t, "ironic")
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.NotContains(t, out, "Alice · since",
		"ironic contributor table should not repeat name in third column")
	assert.Contains(t, out, "since 2025-01-01")
}

// Phase 4 — Bug 1: no color-mix in Chart.js script blocks
func TestBuiltinStyles_noColorMixInScripts(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	scriptPattern := regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`)
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s, err := styles.Load(name)
			require.NoError(t, err)
			checkNoColorMixInScript(t, scriptPattern, s.HTMLTemplateChangelog, name+" changelog")
			checkNoColorMixInScript(t, scriptPattern, s.HTMLTemplateStory, name+" story")
		})
	}
}

func checkNoColorMixInScript(t *testing.T, pattern *regexp.Regexp, tmpl, label string) {
	t.Helper()
	if tmpl == "" {
		return
	}
	matches := pattern.FindAllStringSubmatch(tmpl, -1)
	for _, m := range matches {
		assert.NotContains(t, m[1], "color-mix",
			"%s: color-mix() found in <script> block — Chart.js cannot parse it", label)
	}
}

// Phase 5 — Bug 1: all chart initializations wrapped in setTimeout
func TestBuiltinStyles_chartsWrappedInSetTimeout(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s, err := styles.Load(name)
			require.NoError(t, err)
			checkChartsInSetTimeout(t, s.HTMLTemplateChangelog, name+" changelog")
			checkChartsInSetTimeout(t, s.HTMLTemplateStory, name+" story")
		})
	}
}

func checkChartsInSetTimeout(t *testing.T, tmpl, label string) {
	t.Helper()
	if tmpl == "" {
		return
	}
	lines := strings.Split(tmpl, "\n")
	for i, line := range lines {
		if strings.Contains(line, "new Chart(") {
			start := i - 5
			if start < 0 {
				start = 0
			}
			context := strings.Join(lines[start:i+1], "\n")
			assert.Contains(t, context, "setTimeout",
				"%s line %d: new Chart() must be inside setTimeout", label, i+1)
		}
	}
}

// Phase 5 — Bug 3: ironic story "when people actually worked" not duplicated
func TestRenderStory_ironicNoDuplicateChartTitle(t *testing.T) {
	s := loadTestStyle(t, "ironic")
	out, err := RenderStory("", sampleStoryChronology(), s, FormatHTML)
	require.NoError(t, err)
	count := strings.Count(out, "when people actually worked")
	assert.Equal(t, 1, count,
		"ironic story: chart title should appear exactly once")
}

// Phase 5 — Bug 4: ironic commit icons use .Icon not .Bullet
func TestRenderChangelog_ironicIconsDistinct(t *testing.T) {
	s := loadTestStyle(t, "ironic")
	cl := changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{Type: changelog.TypeBreaking, Commits: []changelog.ParsedCommit{
				{Hash: "aaa0000001", Author: "A", Date: time.Now(),
					Type: changelog.TypeBreaking, Message: "breaking: remove api"},
			}},
			{Type: changelog.TypeFeat, Commits: []changelog.ParsedCommit{
				{Hash: "bbb0000002", Author: "B", Date: time.Now(),
					Type: changelog.TypeFeat, Message: "feat: add thing"},
			}},
		},
	}
	out, err := Render("", cl, s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "!",
		"ironic changelog: breaking commits should show ! icon")
}

// Phase 5 — Bug 7: formal story no inline font-size in stat cells
func TestRenderStory_formalNoInlineFontSizeInStats(t *testing.T) {
	s := loadTestStyle(t, "formal")
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, s, FormatHTML)
	require.NoError(t, err)
	assert.NotContains(t, out, `style="font-size:14px`,
		"formal story: inline font-size overrides should be removed from stat cells")
	assert.NotContains(t, out, `style="font-size:16px`,
		"formal story: inline font-size overrides should be removed from stat cells")
}

// Phase 6 — Content parity tests

func TestAllStyles_changelogHasCommitList(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			out, err := Render("", sampleChangelog(), s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, "add login", "%s changelog: commit message must appear", name)
		})
	}
}

func TestAllStyles_changelogHasTypeChart(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			out, err := Render("", sampleChangelog(), s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, "Chart", "%s changelog: must include a Chart.js chart", name)
		})
	}
}

func TestAllStyles_storyHasTotalCommits(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			ch := sampleStoryChronology()
			out, err := RenderStory("", ch, s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, "62", "%s story: TotalCommits must be visible", name)
		})
	}
}

func TestAllStyles_storyHasFirstDate(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			ch := sampleStoryChronology()
			out, err := RenderStory("", ch, s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, "2025-01-01", "%s story: FirstDate must be visible", name)
		})
	}
}

func TestAllStyles_storyHasMostActiveMonth(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			ch := sampleStoryChronology()
			out, err := RenderStory("", ch, s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, "2025-01", "%s story: most active month must be visible", name)
		})
	}
}

func TestAllStyles_storyHasMilestones(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			ch := sampleStoryChronology()
			out, err := RenderStory("", ch, s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, "v1.0.0", "%s story: tags/milestones must be visible", name)
		})
	}
}

func TestAllStyles_storyHasContributorRankingChart(t *testing.T) {
	styleNames := []string{"formal", "patchnotes", "epic", "ironic"}
	for _, name := range styleNames {
		t.Run(name, func(t *testing.T) {
			s := loadTestStyle(t, name)
			ch := sampleStoryChronology()
			out, err := RenderStory("", ch, s, FormatHTML)
			require.NoError(t, err)
			assert.Contains(t, out, `"Alice"`, "%s story: contributor name must appear in chart data", name)
			assert.Contains(t, out, `"Bob"`, "%s story: contributor name must appear in chart data", name)
		})
	}
}

func TestRepoNameFromPath(t *testing.T) {
	assert.Equal(t, "commitlore", RepoNameFromPath("C:\\Users\\alcil\\MyProjects\\commitlore"))
	assert.Equal(t, "commitlore", RepoNameFromPath("/home/user/commitlore"))
	assert.Equal(t, "repo", RepoNameFromPath("owner/repo"))
	assert.Equal(t, "Repository", RepoNameFromPath("."))
	assert.Equal(t, "Repository", RepoNameFromPath(""))
}
