package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStyleTheme_formal(t *testing.T) {
	s := NewStyleApp()
	theme, err := s.GetStyleTheme("formal")
	require.NoError(t, err)
	assert.Equal(t, "#2563EB", theme.Primary)
	assert.Equal(t, "#F8FAFC", theme.Background)
	assert.Equal(t, "#1E293B", theme.Text)
	assert.Equal(t, "light", theme.Mode)
	assert.Contains(t, theme.FontFamily, "Inter")
}

func TestGetStyleTheme_patchnotes(t *testing.T) {
	s := NewStyleApp()
	theme, err := s.GetStyleTheme("patchnotes")
	require.NoError(t, err)
	assert.Equal(t, "#7C3AED", theme.Primary)
	assert.Equal(t, "#0D0D0D", theme.Background)
	assert.Equal(t, "dark", theme.Mode)
}

func TestGetStyleTheme_missingFields(t *testing.T) {
	s := NewStyleApp()
	theme, err := s.GetStyleTheme("formal")
	require.NoError(t, err)
	assert.NotEmpty(t, theme.Primary)
	assert.NotEmpty(t, theme.Background)
	assert.NotEmpty(t, theme.Surface)
	assert.NotEmpty(t, theme.Text)
	assert.NotEmpty(t, theme.FontFamily)
	assert.NotEmpty(t, theme.Mode)
}

func TestGetStyleTheme_unknownStyle(t *testing.T) {
	s := NewStyleApp()
	_, err := s.GetStyleTheme("nonexistent")
	assert.Error(t, err)
}
