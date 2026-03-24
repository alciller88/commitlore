package llm

import (
	"fmt"
	"net/http"
	"time"
)

const llmTimeout = 30 * time.Second

// New creates a Provider for the given provider name.
// The model parameter overrides the default model for the provider.
// Returns an error if the provider is unknown or if a required API key is missing.
func New(provider, apiKey, baseURL, model string) (Provider, error) {
	resolved, resolvedURL, aliasModel, requiresKey := ResolveAlias(provider, baseURL)
	if requiresKey && apiKey == "" {
		return nil, fmt.Errorf("LLM error: COMMITLORE_LLM_API_KEY is required for provider %q", provider)
	}

	effectiveModel := pickModel(model, aliasModel)
	return newProvider(resolved, apiKey, resolvedURL, effectiveModel)
}

func pickModel(userModel, aliasModel string) string {
	if userModel != "" {
		return userModel
	}
	return aliasModel
}

func newProvider(provider, apiKey, baseURL, model string) (Provider, error) {
	client := &http.Client{Timeout: llmTimeout}

	switch provider {
	case "anthropic":
		return NewAnthropic(apiKey, model, client), nil
	case "openai":
		return NewOpenAI(apiKey, baseURL, model, client), nil
	default:
		return nil, fmt.Errorf("LLM error: unknown provider %q (supported: anthropic, openai, ollama, groq)", provider)
	}
}
