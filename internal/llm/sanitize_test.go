package llm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeRepoData_ShortString(t *testing.T) {
	input := "feat: add login"
	assert.Equal(t, input, SanitizeRepoData(input))
}

func TestSanitizeRepoData_TruncatesAt500(t *testing.T) {
	input := strings.Repeat("a", 600)
	result := SanitizeRepoData(input)
	assert.Len(t, result, 500)
}

func TestSanitizeRepoData_ExactlyAt500(t *testing.T) {
	input := strings.Repeat("b", 500)
	result := SanitizeRepoData(input)
	assert.Len(t, result, 500)
}

func TestSanitizeRepoData_EscapesControlChars(t *testing.T) {
	input := "hello\x00world\x01test\x7f"
	result := SanitizeRepoData(input)
	assert.Equal(t, "helloworldtest", result)
}

func TestSanitizeRepoData_PreservesNewlinesAndTabs(t *testing.T) {
	input := "line1\nline2\ttab"
	result := SanitizeRepoData(input)
	assert.Equal(t, input, result)
}

func TestSanitizeRepoData_EmptyString(t *testing.T) {
	assert.Equal(t, "", SanitizeRepoData(""))
}

func TestSanitizeRepoData_ControlCharsAndTruncate(t *testing.T) {
	input := "\x00" + strings.Repeat("x", 600)
	result := SanitizeRepoData(input)
	assert.Len(t, result, 500)
	assert.NotContains(t, result, "\x00")
}
