package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPrompt_ContainsDelimiters(t *testing.T) {
	result := BuildPrompt("Rewrite this changelog", "some data")
	assert.Contains(t, result, "---DATA START---")
	assert.Contains(t, result, "---DATA END---")
}

func TestBuildPrompt_ContainsPromptAndData(t *testing.T) {
	result := BuildPrompt("My instructions", "commit list here")
	assert.Contains(t, result, "My instructions")
	assert.Contains(t, result, "commit list here")
}

func TestBuildPrompt_DataBetweenDelimiters(t *testing.T) {
	result := BuildPrompt("prompt", "DATA")
	expected := "prompt\n\n---DATA START---\nDATA\n---DATA END---"
	assert.Equal(t, expected, result)
}
