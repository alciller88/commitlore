package narrative

import (
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

func loadStyle(t *testing.T, name string) styles.Style {
	t.Helper()
	s, err := styles.Load(name)
	require.NoError(t, err)
	return s
}

func TestGenerate_formal(t *testing.T) {
	out, err := Generate(sampleChangelog(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "# Changelog")
	assert.Contains(t, out, "add login")
	assert.Contains(t, out, "resolve crash")
}

func TestGenerate_patchnotes(t *testing.T) {
	out, err := Generate(sampleChangelog(), loadStyle(t, "patchnotes"))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "Patch Notes")
}

func TestGenerate_ironic(t *testing.T) {
	out, err := Generate(sampleChangelog(), loadStyle(t, "ironic"))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "What could go wrong")
}

func TestGenerate_epic(t *testing.T) {
	out, err := Generate(sampleChangelog(), loadStyle(t, "epic"))
	require.NoError(t, err)
	assert.NotEmpty(t, out)
	assert.Contains(t, out, "CHRONICLES")
}

func TestGenerate_emptyChangelog(t *testing.T) {
	out, err := Generate(changelog.Changelog{}, loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "# Changelog")
}

func TestGenerate_shortHash(t *testing.T) {
	out, err := Generate(sampleChangelog(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "abc1234567")
}

func TestGenerate_groupLabels(t *testing.T) {
	out, err := Generate(sampleChangelog(), loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "## Features")
	assert.Contains(t, out, "## Bug Fixes")
}

func TestGenerate_usesFixTemplate(t *testing.T) {
	cl := changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{
				Type: changelog.TypeFix,
				Commits: []changelog.ParsedCommit{
					{Hash: "aaa1111111", Author: "Dev", Message: "fix: a bug"},
				},
			},
		},
	}
	out, err := Generate(cl, loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "a bug")
}

func TestGenerate_usesBreakingTemplate(t *testing.T) {
	cl := changelog.Changelog{
		Groups: []changelog.ChangelogGroup{
			{
				Type: changelog.TypeBreaking,
				Commits: []changelog.ParsedCommit{
					{Hash: "bbb2222222", Author: "Dev", Message: "removed old API"},
				},
			},
		},
	}
	out, err := Generate(cl, loadStyle(t, "formal"))
	require.NoError(t, err)
	assert.Contains(t, out, "removed old API")
}
