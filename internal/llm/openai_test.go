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

func TestOpenAIEnrich_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/chat/completions", r.URL.Path)
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		var req openaiRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)
		assert.Equal(t, openaiDefaultModel, req.Model)
		assert.Contains(t, req.Messages[0].Content, "---DATA START---")

		w.Header().Set("Content-Type", "application/json")
		resp := openaiResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{{Message: struct {
				Content string `json:"content"`
			}{Content: "Enriched text"}}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenAI("test-key", server.URL, "", &http.Client{})
	result, err := p.Enrich(context.Background(), "Rewrite", "raw data")
	require.NoError(t, err)
	assert.Equal(t, "Enriched text", result)
}

func TestOpenAIEnrich_CustomBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chat/completions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	p := NewOpenAI("key", server.URL, "", &http.Client{})
	result, err := p.Enrich(context.Background(), "prompt", "data")
	require.NoError(t, err)
	assert.Equal(t, "ok", result)
}

func TestOpenAIEnrich_NoAuthHeaderWithoutKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"local"}}]}`))
	}))
	defer server.Close()

	p := NewOpenAI("", server.URL, "", &http.Client{})
	result, err := p.Enrich(context.Background(), "prompt", "data")
	require.NoError(t, err)
	assert.Equal(t, "local", result)
}

func TestOpenAIEnrich_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
	}))
	defer server.Close()

	p := NewOpenAI("bad-key", server.URL, "", &http.Client{})
	_, err := p.Enrich(context.Background(), "prompt", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authentication failed")
	assert.NotContains(t, err.Error(), "bad-key")
}

func TestOpenAIEnrich_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := NewOpenAI("key", server.URL, "", &http.Client{Timeout: 50 * time.Millisecond})
	_, err := p.Enrich(context.Background(), "prompt", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LLM error")
}

func TestOpenAIEnrich_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer server.Close()

	p := NewOpenAI("key", server.URL, "", &http.Client{})
	_, err := p.Enrich(context.Background(), "prompt", "data")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

func TestOpenAIEnrich_TrailingSlashInURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chat/completions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	p := NewOpenAI("key", server.URL+"/", "", &http.Client{})
	result, err := p.Enrich(context.Background(), "prompt", "data")
	require.NoError(t, err)
	assert.Equal(t, "ok", result)
}
