package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStyleTheme_formal(t *testing.T) {
	s := NewStyleApp()
	theme, err := s.GetStyleTheme("formal")
	require.NoError(t, err)
	assert.Equal(t, "#0969DA", theme.Primary)
	assert.Equal(t, "#FFFFFF", theme.Background)
	assert.Equal(t, "#1A1A2E", theme.Text)
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

func TestGetStyleDetail_formal(t *testing.T) {
	s := NewStyleApp()
	d, err := s.GetStyleDetail("formal")
	require.NoError(t, err)
	assert.Equal(t, "formal", d.Name)
	assert.NotEmpty(t, d.Templates.Header)
	assert.NotEmpty(t, d.Templates.Feature)
	assert.Equal(t, "light", d.Theme.Mode)
	assert.Equal(t, "#0969DA", d.Theme.Colors.Primary)
	assert.Contains(t, d.Theme.Typography.FontFamily, "Inter")
}

func TestIsStyleBuiltIn_formal(t *testing.T) {
	s := NewStyleApp()
	assert.True(t, s.IsStyleBuiltIn("formal"))
}

func TestIsStyleBuiltIn_userStyle(t *testing.T) {
	s := NewStyleApp()
	assert.False(t, s.IsStyleBuiltIn("my-custom-style"))
}

func TestSaveStyleDetail_user(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	dir := t.TempDir()
	os.Setenv("APPDATA", dir)
	os.Setenv("XDG_CONFIG_HOME", dir)

	s := NewStyleApp()
	detail := StyleDetail{
		Name:    "test-save",
		Version: "1.0.0",
		Author:  "tester",
		Templates: TemplatesDetail{
			Header:  "# {{.Version}}",
			Feature: "- {{.Message}}",
		},
		Theme: ThemeDetail{
			Mode:   "dark",
			Colors: ColorsDetail{Primary: "#FF0000"},
		},
	}

	err := s.SaveStyleDetail(detail)
	require.NoError(t, err)

	loaded, err := s.GetStyleDetail("test-save")
	require.NoError(t, err)
	assert.Equal(t, "test-save", loaded.Name)
	assert.Equal(t, "#FF0000", loaded.Theme.Colors.Primary)
	assert.Equal(t, "# {{.Version}}", loaded.Templates.Header)
}

func TestGetStyleTheme_windowControls(t *testing.T) {
	s := NewStyleApp()
	theme, err := s.GetStyleTheme("formal")
	require.NoError(t, err)
	assert.Equal(t, "#888888", theme.WinDefault)
	assert.Equal(t, "#FF5F57", theme.WinClose)
	assert.Equal(t, "#FEBC2E", theme.WinMinimize)
	assert.Equal(t, "#28C840", theme.WinMaximize)
}

func TestGetStyleTheme_windowControls_defaults(t *testing.T) {
	s := NewStyleApp()
	theme, err := s.GetStyleTheme("formal")
	require.NoError(t, err)
	assert.NotEmpty(t, theme.WinDefault)
	assert.NotEmpty(t, theme.WinClose)
	assert.NotEmpty(t, theme.WinMinimize)
	assert.NotEmpty(t, theme.WinMaximize)
}

func TestSaveStyleDetail_builtinRejected(t *testing.T) {
	s := NewStyleApp()
	detail := StyleDetail{Name: "formal", Version: "1.0.0"}
	err := s.SaveStyleDetail(detail)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "built-in")
}
