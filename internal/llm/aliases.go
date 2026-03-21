package llm

// providerAlias maps convenience names to their openai-compatible base URLs.
var providerAliases = map[string]string{
	"ollama": "http://localhost:11434/v1",
	"groq":   "https://api.groq.com/openai/v1",
}

// ResolveAlias checks if the provider name is a convenience alias.
// Returns the resolved provider name, base URL, and whether an API key is required.
func ResolveAlias(provider, baseURL string) (resolvedProvider, resolvedURL string, requiresKey bool) {
	if aliasURL, ok := providerAliases[provider]; ok {
		resolved := "openai"
		if baseURL != "" {
			return resolved, baseURL, provider != "ollama"
		}
		return resolved, aliasURL, provider != "ollama"
	}
	return provider, baseURL, true
}

// IsAlias returns true if the provider name is a known alias.
func IsAlias(provider string) bool {
	_, ok := providerAliases[provider]
	return ok
}
