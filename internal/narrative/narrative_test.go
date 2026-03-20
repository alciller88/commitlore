package narrative

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/alciller88/commitlore/internal/changelog"
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

func TestRender_formalTerminal(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "formal", Format: FormatTerminal})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "add login")
	assert.Contains(t, out, "resolve crash")
}

func TestRender_formalMD(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "formal", Format: FormatMD})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "# Changelog")
	assert.Contains(t, out, "## ")
}

func TestRender_patchnotesTerminal(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "patchnotes", Format: FormatTerminal})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "PATCH NOTES")
}

func TestRender_ironicTerminal(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "ironic", Format: FormatTerminal})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "welcome")
}

func TestRender_epicTerminal(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "epic", Format: FormatTerminal})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "CHRONICLES")
}

func TestRender_jsonValid(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Format: FormatJSON})
	require.NoError(t, err)
	assert.NotEmpty(t, out)

	var parsed jsonChangelog
	err = json.Unmarshal([]byte(out), &parsed)
	require.NoError(t, err)
	assert.Len(t, parsed.Groups, 2)
	assert.Equal(t, "feat", parsed.Groups[0].Type)
	assert.Len(t, parsed.Groups[0].Commits, 2)
}

func TestRender_jsonCommitFields(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Format: FormatJSON})
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

func TestRender_htmlNotImplemented(t *testing.T) {
	_, err := Render(sampleChangelog(), RenderOptions{Format: FormatHTML})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestRender_pdfNotImplemented(t *testing.T) {
	_, err := Render(sampleChangelog(), RenderOptions{Format: FormatPDF})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not implemented")
}

func TestRender_unknownStyleFallsToFormal(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "nonexistent", Format: FormatMD})
	require.NoError(t, err)
	assert.Contains(t, out, "# Changelog")
}

func TestRender_emptyStyleFallsToFormal(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "", Format: FormatMD})
	require.NoError(t, err)
	assert.Contains(t, out, "# Changelog")
}

func TestRender_emptyChangelog(t *testing.T) {
	out, err := Render(changelog.Changelog{}, RenderOptions{Style: "formal", Format: FormatMD})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRender_mdContainsMarkdownSyntax(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "formal", Format: FormatMD})
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(strings.TrimSpace(out), "#"))
	assert.Contains(t, out, "- ")
}

func TestRender_terminalContainsANSI(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "formal", Format: FormatTerminal})
	require.NoError(t, err)
	assert.Contains(t, out, "\033[")
}

func TestRender_shortHash(t *testing.T) {
	out, err := Render(sampleChangelog(), RenderOptions{Style: "formal", Format: FormatMD})
	require.NoError(t, err)
	assert.Contains(t, out, "abc1234")
	assert.NotContains(t, out, "abc1234567")
}
