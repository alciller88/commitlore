package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnthropicEnrich_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "test-key", r.Header.Get("x-api-key"))
		assert.Equal(t, anthropicAPIVersion, r.Header.Get("anthropic-version"))

		var req anthropicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, anthropicDefaultModel, req.Model)
		assert.Len(t, req.Messages, 1)
		assert.Contains(t, req.Messages[0].Content, "---DATA START---")
		assert.Contains(t, req.Messages[0].Content, "---DATA END---")

		w.Header().Set("Content-Type", "application/json")
		resp := anthropicResponse{
			Content: []struct {
				Text string `json:"text"`
			}{{Text: "Enriched changelog text"}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &AnthropicProvider{
		apiKey: "test-key",
		url:    server.URL,
		model:  anthropicDefaultModel,
		client: &http.Client{},
	}

	result, err := p.Enrich(context.Background(), "Rewrite this", "raw data")
	require.NoError(t, err)
	assert.Equal(t, "Enriched changelog text", result)
}

func TestAnthropicEnrich_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer server.Close()

	p := &AnthropicProvider{
		apiKey: "bad-key",
		url:    server.URL,
		model:  anthropicDefaultModel,
		client: &http.Client{},
	}

	_, err := p.Enrich(context.Background(), "prompt", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
	assert.NotContains(t, err.Error(), "bad-key")
}

func TestAnthropicEnrich_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := &AnthropicProvider{
		apiKey: "key",
		url:    server.URL,
		model:  anthropicDefaultModel,
		client: &http.Client{Timeout: 50 * time.Millisecond},
	}

	ctx := context.Background()
	_, err := p.Enrich(ctx, "prompt", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LLM error")
}

func TestAnthropicEnrich_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"content": []}`))
	}))
	defer server.Close()

	p := &AnthropicProvider{
		apiKey: "key",
		url:    server.URL,
		model:  anthropicDefaultModel,
		client: &http.Client{},
	}

	_, err := p.Enrich(context.Background(), "prompt", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

func TestAnthropicEnrich_SanitizesData(t *testing.T) {
	var receivedContent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req anthropicRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		receivedContent = req.Messages[0].Content

		w.Header().Set("Content-Type", "application/json")
		resp := anthropicResponse{
			Content: []struct {
				Text string `json:"text"`
			}{{Text: "ok"}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := &AnthropicProvider{
		apiKey: "key",
		url:    server.URL,
		model:  anthropicDefaultModel,
		client: &http.Client{},
	}

	_, err := p.Enrich(context.Background(), "prompt", "hello\x00world")
	require.NoError(t, err)
	assert.NotContains(t, receivedContent, "\x00")
	assert.Contains(t, receivedContent, "---DATA START---")
}
