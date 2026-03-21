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
func (c *ChangelogApp) Generate(repo, since, until, styleName, llmProvider, llmModel string) (string, error) {
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
	return renderChangelog(cl, styleName, llmProvider, llmModel)
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

	rendered, err := renderer.Render(text, cl, style, renderer.Format("html"))
	if err != nil {
		return "", cleanError(err)
	}

	return rendered, nil
}

func tryEnrich(provider, model, llmPrompt, text string) string {
	if provider == "" || provider == "none" || llmPrompt == "" {
		return text
	}

	apiKey := resolveAPIKey(provider)
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
