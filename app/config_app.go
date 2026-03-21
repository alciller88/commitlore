package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
)

const (
	keyringService = "commitlore"
	maxRecentRepos = 10
)

// RecentRepo represents a recently opened repository.
type RecentRepo struct {
	Path     string `json:"path" yaml:"path"`
	Type     string `json:"type" yaml:"type"`
	LastUsed string `json:"lastUsed" yaml:"last_used"`
}

// LLMConfig holds the LLM provider configuration for the frontend.
type LLMConfig struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	KeyConfigured bool   `json:"keyConfigured"`
}

type appConfig struct {
	RecentRepos []RecentRepo `yaml:"recent_repos"`
	LLM         llmDiskConf  `yaml:"llm"`
}

type llmDiskConf struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

// ConfigApp manages application configuration for the frontend.
type ConfigApp struct{}

// NewConfigApp creates a new ConfigApp instance.
func NewConfigApp() *ConfigApp {
	return &ConfigApp{}
}

// GetRecentRepos returns the list of recent repos from config.
func (c *ConfigApp) GetRecentRepos() ([]RecentRepo, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return cfg.RecentRepos, nil
}

// AddRecentRepo adds a repo to the recent list (max 10, most recent first).
func (c *ConfigApp) AddRecentRepo(path, repoType string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	entry := RecentRepo{Path: path, Type: repoType, LastUsed: now}

	filtered := filterOutRepo(cfg.RecentRepos, path)
	cfg.RecentRepos = append([]RecentRepo{entry}, filtered...)

	if len(cfg.RecentRepos) > maxRecentRepos {
		cfg.RecentRepos = cfg.RecentRepos[:maxRecentRepos]
	}

	return saveConfig(cfg)
}

// SetLLMConfig saves provider and model to config, key to OS keychain.
func (c *ConfigApp) SetLLMConfig(provider, model, key string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	cfg.LLM.Provider = provider
	cfg.LLM.Model = model

	if err := saveConfig(cfg); err != nil {
		return err
	}

	if key != "" {
		return keyring.Set(keyringService, provider, key)
	}
	return nil
}

// GetLLMConfig returns provider, model, and whether a key is configured.
func (c *ConfigApp) GetLLMConfig() (LLMConfig, error) {
	cfg, err := loadConfig()
	if err != nil {
		return LLMConfig{}, err
	}

	configured := false
	if cfg.LLM.Provider != "" {
		_, err := keyring.Get(keyringService, cfg.LLM.Provider)
		configured = err == nil
	}

	return LLMConfig{
		Provider:      cfg.LLM.Provider,
		Model:         cfg.LLM.Model,
		KeyConfigured: configured,
	}, nil
}

// ClearLLMKey removes the API key from the OS keychain.
func (c *ConfigApp) ClearLLMKey(provider string) error {
	err := keyring.Delete(keyringService, provider)
	if err == keyring.ErrNotFound {
		return nil
	}
	return err
}

func filterOutRepo(repos []RecentRepo, path string) []RecentRepo {
	result := make([]RecentRepo, 0, len(repos))
	for _, r := range repos {
		if r.Path != path {
			result = append(result, r)
		}
	}
	return result
}

func configFilePath() (string, error) {
	dir, err := appConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yml"), nil
}

func appConfigDir() (string, error) {
	if runtime.GOOS == "windows" {
		dir := os.Getenv("APPDATA")
		if dir == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		return filepath.Join(dir, "commitlore"), nil
	}
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("getting home dir: %w", err)
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "commitlore"), nil
}

func loadConfig() (appConfig, error) {
	path, err := configFilePath()
	if err != nil {
		return appConfig{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return appConfig{}, nil
		}
		return appConfig{}, fmt.Errorf("reading config: %w", err)
	}

	var cfg appConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return appConfig{}, fmt.Errorf("parsing config: %w", err)
	}

	sort.Slice(cfg.RecentRepos, func(i, j int) bool {
		return cfg.RecentRepos[i].LastUsed > cfg.RecentRepos[j].LastUsed
	})

	return cfg, nil
}

func saveConfig(cfg appConfig) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("serializing config: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}
