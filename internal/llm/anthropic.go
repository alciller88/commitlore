package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	anthropicDefaultURL   = "https://api.anthropic.com/v1/messages"
	anthropicDefaultModel = "claude-haiku-4-5-20251001"
	anthropicAPIVersion   = "2023-06-01"
)

// AnthropicProvider calls the Anthropic Messages API.
type AnthropicProvider struct {
	apiKey string
	url    string
	model  string
	client *http.Client
}

// NewAnthropic creates an Anthropic provider.
func NewAnthropic(apiKey string, httpClient *http.Client) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey: apiKey,
		url:    anthropicDefaultURL,
		model:  anthropicDefaultModel,
		client: httpClient,
	}
}

func (a *AnthropicProvider) Enrich(ctx context.Context, prompt, data string) (string, error) {
	body, err := a.buildRequest(prompt, data)
	if err != nil {
		return "", fmt.Errorf("LLM error: building request: %w", err)
	}
	return a.doRequest(ctx, body)
}

func (a *AnthropicProvider) buildRequest(prompt, data string) ([]byte, error) {
	fullPrompt := BuildPrompt(prompt, SanitizeRepoData(data))
	payload := anthropicRequest{
		Model:     a.model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: fullPrompt},
		},
	}
	return json.Marshal(payload)
}

func (a *AnthropicProvider) doRequest(ctx context.Context, body []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("LLM error: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion)

	return executeAndParse(a.client, req, parseAnthropicResponse)
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func parseAnthropicResponse(data []byte) (string, error) {
	var resp anthropicResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("LLM error: parsing response: %w", err)
	}
	if resp.Error != nil {
		return "", classifyError(resp.Error.Message)
	}
	if len(resp.Content) == 0 {
		return "", fmt.Errorf("LLM error: empty response from provider")
	}
	return resp.Content[0].Text, nil
}
