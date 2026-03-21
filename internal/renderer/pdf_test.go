package renderer

import (
	"testing"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_pdfChangelog(t *testing.T) {
	out, err := Render("", sampleChangelog(), FormatPDF)
	require.NoError(t, err)
	assert.True(t, len(out) > 100, "PDF should have substantial content")
	assert.Contains(t, out, "%PDF", "output should be valid PDF")
	assert.Contains(t, out, "%%EOF", "output should have PDF trailer")
}

func TestRender_pdfChangelogHasFont(t *testing.T) {
	out, err := Render("", sampleChangelog(), FormatPDF)
	require.NoError(t, err)
	assert.Contains(t, out, "/Type /Font")
	assert.Contains(t, out, "Helvetica")
}

func TestRenderStory_pdf(t *testing.T) {
	ch := sampleStoryChronology()
	out, err := RenderStory("", ch, FormatPDF)
	require.NoError(t, err)
	assert.True(t, len(out) > 100, "PDF should have substantial content")
	assert.Contains(t, out, "%PDF")
	assert.Contains(t, out, "%%EOF")
}

func TestRenderStory_pdfLargerWithContent(t *testing.T) {
	empty, err := RenderStory("", git.Chronology{}, FormatPDF)
	require.NoError(t, err)

	full, err := RenderStory("", sampleStoryChronology(), FormatPDF)
	require.NoError(t, err)

	assert.Greater(t, len(full), len(empty), "PDF with content should be larger than empty")
}

func TestRenderStory_pdfEmptyChronology(t *testing.T) {
	out, err := RenderStory("", git.Chronology{}, FormatPDF)
	require.NoError(t, err)
	assert.True(t, len(out) > 0)
	assert.Contains(t, out, "%PDF")
}
