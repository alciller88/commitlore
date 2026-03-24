package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Anthropic(t *testing.T) {
	p, err := New("anthropic", "sk-test", "", "")
	require.NoError(t, err)
	assert.IsType(t, &AnthropicProvider{}, p)
}

func TestNew_AnthropicWithModel(t *testing.T) {
	p, err := New("anthropic", "sk-test", "", "claude-sonnet-4-6")
	require.NoError(t, err)
	ap, ok := p.(*AnthropicProvider)
	require.True(t, ok)
	assert.Equal(t, "claude-sonnet-4-6", ap.model)
}

func TestNew_OpenAI(t *testing.T) {
	p, err := New("openai", "sk-test", "", "")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, openaiDefaultModel, oai.model)
}

func TestNew_OpenAIWithBaseURL(t *testing.T) {
	p, err := New("openai", "sk-test", "http://custom:8080/v1", "")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, "http://custom:8080/v1", oai.baseURL)
}

func TestNew_OllamaNoKeyRequired(t *testing.T) {
	p, err := New("ollama", "", "", "")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, "llama3.2", oai.model)
}

func TestNew_OllamaWithCustomModel(t *testing.T) {
	p, err := New("ollama", "", "", "mistral")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, "mistral", oai.model)
}

func TestNew_GroqRequiresKey(t *testing.T) {
	_, err := New("groq", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMMITLORE_LLM_API_KEY is required")
}

func TestNew_GroqWithKey(t *testing.T) {
	p, err := New("groq", "gsk-test", "", "")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, "llama-3.3-70b-versatile", oai.model)
}

func TestNew_GroqWithCustomModel(t *testing.T) {
	p, err := New("groq", "gsk-test", "", "mixtral-8x7b-32768")
	require.NoError(t, err)
	oai, ok := p.(*OpenAIProvider)
	require.True(t, ok)
	assert.Equal(t, "mixtral-8x7b-32768", oai.model)
}

func TestNew_UnknownProvider(t *testing.T) {
	_, err := New("gemini", "key", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown provider")
}

func TestNew_AnthropicMissingKey(t *testing.T) {
	_, err := New("anthropic", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMMITLORE_LLM_API_KEY is required")
}

func TestNew_OpenAIMissingKey(t *testing.T) {
	_, err := New("openai", "", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "COMMITLORE_LLM_API_KEY is required")
}
