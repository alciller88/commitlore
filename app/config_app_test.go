package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func setupTestConfig(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	origAppData := os.Getenv("APPDATA")
	origXDG := os.Getenv("XDG_CONFIG_HOME")

	os.Setenv("APPDATA", dir)
	os.Setenv("XDG_CONFIG_HOME", dir)

	return dir, func() {
		os.Setenv("APPDATA", origAppData)
		os.Setenv("XDG_CONFIG_HOME", origXDG)
	}
}

func TestAddRecentRepo_addsEntry(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	err := c.AddRecentRepo("/path/to/repo", "local")
	require.NoError(t, err)

	repos, err := c.GetRecentRepos()
	require.NoError(t, err)
	assert.Len(t, repos, 1)
	assert.Equal(t, "/path/to/repo", repos[0].Path)
	assert.Equal(t, "local", repos[0].Type)
}

func TestAddRecentRepo_mostRecentFirst(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	require.NoError(t, c.AddRecentRepo("/first", "local"))
	time.Sleep(10 * time.Millisecond)
	require.NoError(t, c.AddRecentRepo("/second", "github"))

	repos, err := c.GetRecentRepos()
	require.NoError(t, err)
	assert.Len(t, repos, 2)
	assert.Equal(t, "/second", repos[0].Path)
	assert.Equal(t, "/first", repos[1].Path)
}

func TestAddRecentRepo_deduplicates(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	require.NoError(t, c.AddRecentRepo("/repo", "local"))
	require.NoError(t, c.AddRecentRepo("/other", "github"))
	require.NoError(t, c.AddRecentRepo("/repo", "local"))

	repos, err := c.GetRecentRepos()
	require.NoError(t, err)
	assert.Len(t, repos, 2)
	assert.Equal(t, "/repo", repos[0].Path)
}

func TestAddRecentRepo_maxTenEntries(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	for i := 0; i < 15; i++ {
		require.NoError(t, c.AddRecentRepo(
			filepath.Join("/repo", string(rune('a'+i))),
			"local",
		))
	}

	repos, err := c.GetRecentRepos()
	require.NoError(t, err)
	assert.Len(t, repos, maxRecentRepos)
}

func TestGetRecentRepos_emptyConfig(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	repos, err := c.GetRecentRepos()
	require.NoError(t, err)
	assert.Empty(t, repos)
}

func TestGetLLMConfig_neverReturnsKey(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	cfg, err := c.GetLLMConfig()
	require.NoError(t, err)
	assert.Equal(t, "", cfg.Provider)
	assert.Equal(t, "", cfg.Model)
	assert.False(t, cfg.KeyConfigured)
}

func TestSetLLMConfig_savesProviderAndModel(t *testing.T) {
	dir, cleanup := setupTestConfig(t)
	defer cleanup()

	c := NewConfigApp()
	err := c.SetLLMConfig("anthropic", "claude-haiku-4-5-20251001", "")
	require.NoError(t, err)

	path := filepath.Join(dir, "commitlore", "config.yml")
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var cfg appConfig
	require.NoError(t, yaml.Unmarshal(data, &cfg))
	assert.Equal(t, "anthropic", cfg.LLM.Provider)
	assert.Equal(t, "claude-haiku-4-5-20251001", cfg.LLM.Model)

	// Verify no key in file
	assert.NotContains(t, string(data), "sk-")
}

func TestConfigPersistence_surviveReload(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	c1 := NewConfigApp()
	require.NoError(t, c1.AddRecentRepo("/persistent", "local"))
	require.NoError(t, c1.SetLLMConfig("openai", "gpt-4o-mini", ""))

	c2 := NewConfigApp()
	repos, err := c2.GetRecentRepos()
	require.NoError(t, err)
	assert.Len(t, repos, 1)
	assert.Equal(t, "/persistent", repos[0].Path)

	cfg, err := c2.GetLLMConfig()
	require.NoError(t, err)
	assert.Equal(t, "openai", cfg.Provider)
	assert.Equal(t, "gpt-4o-mini", cfg.Model)
}
