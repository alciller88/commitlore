package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveAlias_Ollama(t *testing.T) {
	provider, url, requiresKey := ResolveAlias("ollama", "")
	assert.Equal(t, "openai", provider)
	assert.Equal(t, "http://localhost:11434/v1", url)
	assert.False(t, requiresKey)
}

func TestResolveAlias_Groq(t *testing.T) {
	provider, url, requiresKey := ResolveAlias("groq", "")
	assert.Equal(t, "openai", provider)
	assert.Equal(t, "https://api.groq.com/openai/v1", url)
	assert.True(t, requiresKey)
}

func TestResolveAlias_OllamaWithCustomURL(t *testing.T) {
	provider, url, _ := ResolveAlias("ollama", "http://custom:8080/v1")
	assert.Equal(t, "openai", provider)
	assert.Equal(t, "http://custom:8080/v1", url)
}

func TestResolveAlias_NotAnAlias(t *testing.T) {
	provider, url, requiresKey := ResolveAlias("anthropic", "")
	assert.Equal(t, "anthropic", provider)
	assert.Equal(t, "", url)
	assert.True(t, requiresKey)
}

func TestResolveAlias_OpenAIPassthrough(t *testing.T) {
	provider, url, requiresKey := ResolveAlias("openai", "https://custom.api.com/v1")
	assert.Equal(t, "openai", provider)
	assert.Equal(t, "https://custom.api.com/v1", url)
	assert.True(t, requiresKey)
}

func TestIsAlias_KnownAliases(t *testing.T) {
	assert.True(t, IsAlias("ollama"))
	assert.True(t, IsAlias("groq"))
}

func TestIsAlias_NotAliases(t *testing.T) {
	assert.False(t, IsAlias("anthropic"))
	assert.False(t, IsAlias("openai"))
	assert.False(t, IsAlias("unknown"))
}
