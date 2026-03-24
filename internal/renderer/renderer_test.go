package renderer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/styles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleChangelog() changelog.Changelog {
	base := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	return changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{
				Type: changelog.TypeFeat,
				Commits: []changelog.ParsedCommit{
					{Hash: "abc1234567", Author: "Alice", Email: "a@x.com", Date: base, Type: changelog.TypeFeat, Message: "feat: add login"},
					{Hash: "def7890123", Author: "Bob", Email: "b@x.com", Date: base, Type: changelog.TypeFeat, Message: "feat: add logout"},
				},
			},
			{
				Type: changelog.TypeFix,
				Commits: []changelog.ParsedCommit{
					{Hash: "ghi4567890", Author: "Alice", Email: "a@x.com", Date: base, Type: changelog.TypeFix, Message: "fix: resolve crash"},
				},
			},
		},
	}
}

func loadTestStyle(t *testing.T, name string) styles.Style {
	t.Helper()
	s, err := styles.Load(name)
	require.NoError(t, err)
	return s
}

func TestRender_terminalAddsANSI(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out, err := Render("# Changelog\n## Features\n- item", changelog.Changelog{}, s, FormatTerminal)
	require.NoError(t, err)
	assert.Contains(t, out, "\033[")
}

func TestRender_mdPassesThrough(t *testing.T) {
	content := "# Changelog\n## Features\n- item"
	out, err := Render(content, changelog.Changelog{}, styles.Style{}, FormatMD)
	require.NoError(t, err)
	assert.Equal(t, content, out)
}

func TestRender_jsonValid(t *testing.T) {
	out, err := Render("", sampleChangelog(), styles.Style{}, FormatJSON)
	require.NoError(t, err)

	var parsed jsonChangelog
	err = json.Unmarshal([]byte(out), &parsed)
	require.NoError(t, err)
	assert.Len(t, parsed.Groups, 2)
	assert.Equal(t, "feat", parsed.Groups[0].Type)
	assert.Len(t, parsed.Groups[0].Commits, 2)
}

func TestRender_jsonCommitFields(t *testing.T) {
	out, err := Render("", sampleChangelog(), styles.Style{}, FormatJSON)
	require.NoError(t, err)

	var parsed jsonChangelog
	err = json.Unmarshal([]byte(out), &parsed)
	require.NoError(t, err)

	c := parsed.Groups[0].Commits[0]
	assert.Equal(t, "abc1234567", c.Hash)
	assert.Equal(t, "Alice", c.Author)
	assert.Equal(t, "a@x.com", c.Email)
	assert.Equal(t, "2025-06-01", c.Date)
	assert.Equal(t, "feat", c.Type)
}

func TestRender_htmlReturnsHTML(t *testing.T) {
	s := loadTestStyle(t, "formal")
	out, err := Render("", sampleChangelog(), s, FormatHTML)
	require.NoError(t, err)
	assert.Contains(t, out, "<html")
}

func TestRender_pdfReturnsError(t *testing.T) {
	_, err := Render("", sampleChangelog(), styles.Style{}, FormatPDF)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "removed")
}

func TestRender_unknownFormatPassesThrough(t *testing.T) {
	content := "some text"
	out, err := Render(content, changelog.Changelog{}, styles.Style{}, Format("unknown"))
	require.NoError(t, err)
	assert.Equal(t, content, out)
}
