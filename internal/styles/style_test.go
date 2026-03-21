package styles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_formal(t *testing.T) {
	s, err := Load("formal")
	require.NoError(t, err)
	assert.Equal(t, "formal", s.Name)
	assert.Equal(t, "1.0.0", s.Version)
	assert.NotEmpty(t, s.Description)
	assert.NotEmpty(t, s.Templates.Header)
	assert.NotEmpty(t, s.Templates.Feature)
}

func TestLoad_patchnotes(t *testing.T) {
	s, err := Load("patchnotes")
	require.NoError(t, err)
	assert.Equal(t, "patchnotes", s.Name)
	assert.NotEmpty(t, s.Templates.Header)
	assert.Contains(t, s.Templates.Header, "PATCH NOTES")
}

func TestLoad_ironic(t *testing.T) {
	s, err := Load("ironic")
	require.NoError(t, err)
	assert.Equal(t, "ironic", s.Name)
	assert.Contains(t, s.Templates.Header, "Somehow it works")
}

func TestLoad_epic(t *testing.T) {
	s, err := Load("epic")
	require.NoError(t, err)
	assert.Equal(t, "epic", s.Name)
	assert.Contains(t, s.Templates.Header, "CHRONICLES")
}

func TestLoad_emptyNameDefaultsToFormal(t *testing.T) {
	s, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, "formal", s.Name)
}

func TestLoad_unknownStyle(t *testing.T) {
	_, err := Load("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestLoad_allBuiltinsHaveLLMPrompt(t *testing.T) {
	for _, name := range ListBuiltin() {
		s, err := Load(name)
		require.NoError(t, err, "style: %s", name)
		assert.NotEmpty(t, s.LLMPrompt, "style %s should have llm_prompt", name)
	}
}

func TestListBuiltin(t *testing.T) {
	names := ListBuiltin()
	assert.Len(t, names, 4)
	assert.Contains(t, names, "formal")
	assert.Contains(t, names, "patchnotes")
	assert.Contains(t, names, "ironic")
	assert.Contains(t, names, "epic")
}

func TestIsBuiltin(t *testing.T) {
	assert.True(t, IsBuiltin("formal"))
	assert.True(t, IsBuiltin("epic"))
	assert.False(t, IsBuiltin("custom"))
	assert.False(t, IsBuiltin(""))
}

func TestLoad_themeLoaded(t *testing.T) {
	s, err := Load("epic")
	require.NoError(t, err)
	assert.Equal(t, "dark", s.Theme.Mode)
	assert.Equal(t, "#D4AF37", s.Theme.Colors.Primary)
	assert.Equal(t, "glassmorphism", s.Theme.CardStyle)
	assert.True(t, s.Theme.Animations)
}

func TestLoad_terminalLoaded(t *testing.T) {
	s, err := Load("patchnotes")
	require.NoError(t, err)
	assert.Equal(t, "magenta", s.Terminal.Colors.Header)
	assert.Equal(t, "normal", s.Terminal.Density)
	assert.NotEmpty(t, s.Terminal.Decorators.Separator)
}

func TestLoad_vocabularyLoaded(t *testing.T) {
	s, err := Load("epic")
	require.NoError(t, err)
	assert.Equal(t, "dark curse", s.Vocabulary["bug"])
	assert.Equal(t, "vanquished", s.Vocabulary["fix"])
}

func TestValidate_invalidCardStyle(t *testing.T) {
	s := Style{Name: "test", Templates: Templates{Header: "h"}, Theme: Theme{CardStyle: "invalid"}}
	err := validate(s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "card_style")
}

func TestValidate_invalidDensity(t *testing.T) {
	s := Style{Name: "test", Templates: Templates{Header: "h"}, Terminal: Terminal{Density: "invalid"}}
	err := validate(s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "density")
}

func TestValidate_invalidMode(t *testing.T) {
	s := Style{Name: "test", Templates: Templates{Header: "h"}, Theme: Theme{Mode: "invalid"}}
	err := validate(s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mode")
}

func TestValidate_emptyOptionalFieldsPass(t *testing.T) {
	s := Style{Name: "test", Templates: Templates{Header: "h"}}
	err := validate(s)
	assert.NoError(t, err)
}

func TestValidate_validValues(t *testing.T) {
	s := Style{
		Name:      "test",
		Templates: Templates{Header: "h"},
		Theme:     Theme{Mode: "dark", CardStyle: "glassmorphism"},
		Terminal:  Terminal{Density: "compact"},
	}
	err := validate(s)
	assert.NoError(t, err)
}
