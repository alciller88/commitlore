package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validShipstyle = `name: "test-community"
version: "1.0.0"
description: "A community style"
author: "tester"
templates:
  header: "# Test"
  feature: "- {{.Message}}"
`

const validCatalog = `[
  {
    "name": "test-community",
    "description": "A community style",
    "author": "tester",
    "version": "1.0.0",
    "tags": ["dark"],
    "preview": "PREVIEW_URL",
    "download": "DOWNLOAD_URL"
  }
]`

func setTestEnv(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("APPDATA", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
}

func newTestMarketplace(client *http.Client) *MarketplaceApp {
	return &MarketplaceApp{httpClient: client}
}

func TestFetchCatalog_parsesCatalog(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(validCatalog))
	}))
	defer srv.Close()

	m := newTestMarketplace(srv.Client())
	entries, err := m.fetchCatalogFromURL(srv.URL)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "test-community", entries[0].Name)
	assert.Equal(t, "tester", entries[0].Author)
	assert.Equal(t, []string{"dark"}, entries[0].Tags)
}

func TestFetchCatalog_networkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	m := newTestMarketplace(srv.Client())
	_, err := m.fetchCatalogFromURL(srv.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}

func TestInstallStyle_savesFile(t *testing.T) {
	setTestEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(validShipstyle))
	}))
	defer srv.Close()

	m := newTestMarketplace(srv.Client())
	err := m.InstallStyle(srv.URL+"/style.shipstyle", "test-community")
	require.NoError(t, err)

	assert.True(t, m.IsInstalled("test-community"))
}

func TestInstallStyle_rejectsBuiltinName(t *testing.T) {
	setTestEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(validShipstyle))
	}))
	defer srv.Close()

	m := newTestMarketplace(srv.Client())
	err := m.InstallStyle(srv.URL, "formal")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "built-in name")
}

func TestInstallStyle_rejectsInvalidSchema(t *testing.T) {
	setTestEnv(t)

	invalidStyle := `name: "bad"
version: "1.0.0"
unknown_field: "should fail"
templates:
  header: "# Test"
  feature: "- {{.Message}}"
`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(invalidStyle))
	}))
	defer srv.Close()

	m := newTestMarketplace(srv.Client())
	err := m.InstallStyle(srv.URL, "bad-style")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid style")
}

const validShipstyleES = `name: "test-community"
language: "es"
version: "1.0.0"
description: "Un estilo comunitario"
author: "tester"
templates:
  header: "# Prueba"
  feature: "- {{.Message}}"
`

func TestInstallStyleWithVariants_savesBaseAndVariant(t *testing.T) {
	setTestEnv(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/base.shipstyle", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(validShipstyle))
	})
	mux.HandleFunc("/es.shipstyle", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(validShipstyleES))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	m := newTestMarketplace(srv.Client())
	variants := map[string]string{"es": srv.URL + "/es.shipstyle"}
	err := m.InstallStyleWithVariants(srv.URL+"/base.shipstyle", "test-community", variants)
	require.NoError(t, err)

	assert.True(t, m.IsInstalled("test-community"))

	dir := os.Getenv("APPDATA")
	if dir == "" {
		dir = os.Getenv("XDG_CONFIG_HOME")
	}
	esPath := filepath.Join(dir, "commitlore", "styles", "test-community.es.shipstyle")
	_, err = os.Stat(esPath)
	assert.NoError(t, err, "Spanish variant should exist")
}

func TestIsInstalled_true(t *testing.T) {
	setTestEnv(t)

	m := NewMarketplaceApp()

	dir := os.Getenv("APPDATA")
	if dir == "" {
		dir = os.Getenv("XDG_CONFIG_HOME")
	}
	stylesDir := filepath.Join(dir, "commitlore", "styles")
	require.NoError(t, os.MkdirAll(stylesDir, 0750))
	require.NoError(t, os.WriteFile(
		filepath.Join(stylesDir, "existing.shipstyle"),
		[]byte(validShipstyle), 0644,
	))

	assert.True(t, m.IsInstalled("existing"))
}

func TestIsInstalled_false(t *testing.T) {
	setTestEnv(t)

	m := NewMarketplaceApp()
	assert.False(t, m.IsInstalled("nonexistent"))
}
