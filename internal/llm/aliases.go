package llm

// aliasConfig holds the resolved settings for a provider alias.
type aliasConfig struct {
	baseURL     string
	model       string
	requiresKey bool
}

// providerAliases maps convenience names to their openai-compatible config.
var providerAliases = map[string]aliasConfig{
	"ollama": {
		baseURL:     "http://localhost:11434/v1",
		model:       "llama3.2",
		requiresKey: false,
	},
	"groq": {
		baseURL:     "https://api.groq.com/openai/v1",
		model:       "llama-3.3-70b-versatile",
		requiresKey: true,
	},
}

// ResolveAlias checks if the provider name is a convenience alias.
// Returns the resolved provider name, base URL, default model, and whether
// an API key is required.
func ResolveAlias(provider, baseURL string) (resolvedProvider, resolvedURL, defaultModel string, requiresKey bool) {
	cfg, ok := providerAliases[provider]
	if !ok {
		return provider, baseURL, "", true
	}

	resolved := "openai"
	if baseURL != "" {
		return resolved, baseURL, cfg.model, cfg.requiresKey
	}
	return resolved, cfg.baseURL, cfg.model, cfg.requiresKey
}

// IsAlias returns true if the provider name is a known alias.
func IsAlias(provider string) bool {
	_, ok := providerAliases[provider]
	return ok
}
