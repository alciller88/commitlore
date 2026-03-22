package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/llm"
	"github.com/alciller88/commitlore/internal/narrative"
	"github.com/alciller88/commitlore/internal/renderer"
	"github.com/alciller88/commitlore/internal/styles"
	"github.com/zalando/go-keyring"
)

const enrichTimeout = 30 * time.Second

// ChangelogApp exposes changelog generation to the frontend.
type ChangelogApp struct{}

// NewChangelogApp creates a new ChangelogApp instance.
func NewChangelogApp() *ChangelogApp {
	return &ChangelogApp{}
}

// Generate produces a changelog as HTML from the given repo.
// LLM config is read automatically from Settings (config.yml + keychain).
func (c *ChangelogApp) Generate(repo, since, until, styleName string) (string, error) {
	opts, err := buildOpts("", since, until, 0)
	if err != nil {
		return "", err
	}

	commits, err := fetchCommits(repo, opts)
	if err != nil {
		return "", cleanError(err)
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found matching the given filters")
	}

	cl := changelog.GroupCommits(commits)
	provider, model := loadLLMSettings()
	return renderChangelog(cl, styleName, provider, model)
}

func renderChangelog(cl changelog.Changelog, styleName, llmProvider, llmModel string) (string, error) {
	style, err := styles.Load(styleName)
	if err != nil {
		return "", cleanError(err)
	}

	text, err := narrative.Generate(cl, style)
	if err != nil {
		return "", cleanError(err)
	}

	text = tryEnrich(llmProvider, llmModel, style.LLMPrompt, text)

	override := buildHTMLThemeOverride(styleName)
	rendered, err := renderer.RenderWithTheme(text, cl, style, renderer.Format("html"), override)
	if err != nil {
		return "", cleanError(err)
	}

	return rendered, nil
}

func buildHTMLThemeOverride(styleName string) *renderer.HTMLTheme {
	st, err := styles.Load(styleName)
	if err != nil {
		return nil
	}
	t := st.Theme
	return &renderer.HTMLTheme{
		Background: t.Colors.Background,
		Surface:    t.Colors.Surface,
		Text:       t.Colors.Text,
		Primary:    t.Colors.Primary,
		Secondary:  t.Colors.Secondary,
		Accent:     t.Colors.Accent,
		Border:     t.Colors.Border,
		FontFamily: t.Typography.FontFamily,
		Mode:       t.Mode,
	}
}

func tryEnrich(provider, model, llmPrompt, text string) string {
	if provider == "" || llmPrompt == "" {
		return text
	}

	apiKey := resolveAPIKey(provider)
	if apiKey == "" {
		return text
	}

	baseURL := os.Getenv("COMMITLORE_LLM_BASE_URL")

	p, err := llm.New(provider, apiKey, baseURL, model)
	if err != nil {
		return text
	}

	ctx, cancel := context.WithTimeout(context.Background(), enrichTimeout)
	defer cancel()

	result, err := p.Enrich(ctx, llmPrompt, text)
	if err != nil {
		return text
	}

	return result
}

// resolveAPIKey checks env var first, then OS keychain.
func resolveAPIKey(provider string) string {
	if key := os.Getenv("COMMITLORE_LLM_API_KEY"); key != "" {
		return key
	}

	key, err := keyring.Get(keyringService, provider)
	if err != nil {
		return ""
	}
	return key
}

// loadLLMSettings reads provider and model from config.yml.
// Returns empty strings if no LLM is configured.
func loadLLMSettings() (string, string) {
	cfg, err := loadConfig()
	if err != nil {
		return "", ""
	}
	return cfg.LLM.Provider, cfg.LLM.Model
}
