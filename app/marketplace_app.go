package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alciller88/commitlore/internal/styles"
)

// CatalogURL is the URL of the official style catalog index.
const CatalogURL = "https://raw.githubusercontent.com/alciller88/commitlore-styles/main/index.json"

const maxCatalogSize = 5 * 1024 * 1024 // 5MB

// MarketplaceEntry represents a single style in the catalog.
type MarketplaceEntry struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	Preview     string   `json:"preview"`
	Download    string   `json:"download"`
}

// MarketplaceApp exposes marketplace operations to the frontend.
type MarketplaceApp struct {
	httpClient *http.Client
}

// NewMarketplaceApp creates a new MarketplaceApp with the default HTTP client.
func NewMarketplaceApp() *MarketplaceApp {
	return &MarketplaceApp{httpClient: http.DefaultClient}
}

// FetchCatalog downloads and parses the style catalog from the official repository.
func (m *MarketplaceApp) FetchCatalog() ([]MarketplaceEntry, error) {
	return m.fetchCatalogFromURL(CatalogURL)
}

func (m *MarketplaceApp) fetchCatalogFromURL(url string) ([]MarketplaceEntry, error) {
	resp, err := m.httpClient.Get(url) //nolint:noctx // simple GET, no context needed
	if err != nil {
		return nil, fmt.Errorf("fetching catalog: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching catalog: HTTP %d", resp.StatusCode)
	}

	data, err := readCatalogBody(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries []MarketplaceEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing catalog: %w", err)
	}
	return entries, nil
}

func readCatalogBody(r io.Reader) ([]byte, error) {
	limited := io.LimitReader(r, int64(maxCatalogSize)+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading catalog: %w", err)
	}
	if len(data) > maxCatalogSize {
		return nil, fmt.Errorf("catalog exceeds maximum size of 5MB")
	}
	return data, nil
}

// InstallStyle downloads a .shipstyle from the given URL, validates it,
// and saves it to the user styles directory.
func (m *MarketplaceApp) InstallStyle(downloadURL, name string) error {
	if styles.IsBuiltin(name) {
		return fmt.Errorf("cannot install style with built-in name %q", name)
	}
	if err := styles.ValidateName(name); err != nil {
		return err
	}

	data, err := m.downloadFile(downloadURL)
	if err != nil {
		return fmt.Errorf("downloading style: %w", err)
	}

	s, err := styles.ParseStyleStrict(data)
	if err != nil {
		return fmt.Errorf("invalid style: %w", err)
	}

	if s.LLMPrompt != "" {
		fmt.Printf("Warning: style %q contains an llm_prompt field. "+
			"Review it before using with --llm.\n", name)
	}

	return saveStyleData(name, data)
}

func (m *MarketplaceApp) downloadFile(url string) ([]byte, error) {
	resp, err := m.httpClient.Get(url) //nolint:noctx // simple GET
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	limited := io.LimitReader(resp.Body, 1024*1024+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if len(data) > 1024*1024 {
		return nil, fmt.Errorf("style file exceeds maximum size of 1MB")
	}
	return data, nil
}

func saveStyleData(name string, data []byte) error {
	dir, err := styles.UserStylesDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name+".shipstyle")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("saving style: %w", err)
	}
	return nil
}

// IsInstalled returns true if a style with the given name exists
// in the user styles directory.
func (m *MarketplaceApp) IsInstalled(name string) bool {
	dir, err := styles.UserStylesDir()
	if err != nil {
		return false
	}
	path := filepath.Join(dir, name+".shipstyle")
	_, err = os.Stat(path)
	return err == nil
}
