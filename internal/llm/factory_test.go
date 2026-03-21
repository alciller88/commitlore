package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Anthropic(t *testing.T) {
	p, err := New("anthropic", "sk-test", "")
	require.NoError(t, err)
	assert.IsType(t, &AnthropicProvider{}, p)
}

func TestNew_OpenAI(t *testing.T) {
	p, err := New("openai", "sk-test", "")
	require.NoError(t, err)
	assert.IsType(t, &OpenAIProvider{}, p)
}

func TestNew_OpenAIWithBaseURL(t *testing.T) {
	p, err := New("openai", "sk-test", "http://custom:8080/v1")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, "http://custom:8080/v1", oai.baseURL)
}

func TestNew_OllamaNoKeyRequired(t *testing.T) {
	p, err := New("ollama", "", "")
	require.NoError(t, err)
	assert.IsType(t, &OpenAIProvider{}, p)
}

func TestNew_GroqRequiresKey(t *testing.T) {
	_, err := New("groq", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMMITLORE_LLM_API_KEY is required")
}

func TestNew_GroqWithKey(t *testing.T) {
	p, err := New("groq", "gsk-test", "")
	require.NoError(t, err)
	assert.IsType(t, &OpenAIProvider{}, p)
}

func TestNew_UnknownProvider(t *testing.T) {
	_, err := New("gemini", "key", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown provider")
}

func TestNew_AnthropicMissingKey(t *testing.T) {
	_, err := New("anthropic", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMMITLORE_LLM_API_KEY is required")
}

func TestNew_OpenAIMissingKey(t *testing.T) {
	_, err := New("openai", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMMITLORE_LLM_API_KEY is required")
}
