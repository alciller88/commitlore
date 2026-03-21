package narrative

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyVocabulary_basicReplacement(t *testing.T) {
	vocab := map[string]string{"bug": "heresy", "fix": "purge"}
	result := ApplyVocabulary("Fixed a bug in the code", vocab)
	assert.Contains(t, result, "heresy")
	assert.Contains(t, result, "purge")
	assert.NotContains(t, result, "bug")
}

func TestApplyVocabulary_caseInsensitive(t *testing.T) {
	vocab := map[string]string{"bug": "heresy"}
	result := ApplyVocabulary("Found a BUG and a Bug", vocab)
	assert.Equal(t, "Found a heresy and a heresy", result)
}

func TestApplyVocabulary_emptyVocab(t *testing.T) {
	result := ApplyVocabulary("no changes here", nil)
	assert.Equal(t, "no changes here", result)
}

func TestApplyVocabulary_noMatch(t *testing.T) {
	vocab := map[string]string{"bug": "heresy"}
	result := ApplyVocabulary("all good", vocab)
	assert.Equal(t, "all good", result)
}

func TestApplyVocabulary_multipleOccurrences(t *testing.T) {
	vocab := map[string]string{"commit": "deed"}
	result := ApplyVocabulary("commit after commit", vocab)
	assert.Equal(t, "deed after deed", result)
}
