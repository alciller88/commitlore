package llm

import (
	"fmt"
	"net/http"
	"time"
)

const llmTimeout = 30 * time.Second

// New creates a Provider for the given provider name.
// Returns an error if the provider is unknown or if a required API key is missing.
func New(provider, apiKey, baseURL string) (Provider, error) {
	resolved, resolvedURL, requiresKey := ResolveAlias(provider, baseURL)
	if requiresKey && apiKey == "" {
		return nil, fmt.Errorf("LLM error: COMMITLORE_LLM_API_KEY is required for provider %q", provider)
	}
	return newProvider(resolved, apiKey, resolvedURL)
}

func newProvider(provider, apiKey, baseURL string) (Provider, error) {
	client := &http.Client{Timeout: llmTimeout}

	switch provider {
	case "anthropic":
		return NewAnthropic(apiKey, client), nil
	case "openai":
		return NewOpenAI(apiKey, baseURL, client), nil
	default:
		return nil, fmt.Errorf("LLM error: unknown provider %q (supported: anthropic, openai, ollama, groq)", provider)
	}
}
