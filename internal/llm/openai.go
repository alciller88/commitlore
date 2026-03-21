package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	openaiDefaultURL   = "https://api.openai.com/v1"
	openaiDefaultModel = "gpt-4o-mini"
)

// OpenAIProvider calls the OpenAI-compatible Chat Completions API.
type OpenAIProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewOpenAI creates an OpenAI-compatible provider.
func NewOpenAI(apiKey, baseURL, model string, httpClient *http.Client) *OpenAIProvider {
	if baseURL == "" {
		baseURL = openaiDefaultURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	if model == "" {
		model = openaiDefaultModel
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  httpClient,
	}
}

func (o *OpenAIProvider) Enrich(ctx context.Context, prompt, data string) (string, error) {
	body, err := o.buildRequest(prompt, data)
	if err != nil {
		return "", fmt.Errorf("LLM error: building request: %w", err)
	}
	return o.doRequest(ctx, body)
}

func (o *OpenAIProvider) buildRequest(prompt, data string) ([]byte, error) {
	fullPrompt := BuildPrompt(prompt, SanitizeRepoData(data))
	payload := openaiRequest{
		Model: o.model,
		Messages: []openaiMessage{
			{Role: "user", Content: fullPrompt},
		},
	}
	return json.Marshal(payload)
}

func (o *OpenAIProvider) doRequest(ctx context.Context, body []byte) (string, error) {
	url := o.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("LLM error: creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if o.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+o.apiKey)
	}

	return executeAndParse(o.client, req, parseOpenAIResponse)
}

type openaiRequest struct {
	Model    string          `json:"model"`
	Messages []openaiMessage `json:"messages"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func parseOpenAIResponse(data []byte) (string, error) {
	var resp openaiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("LLM error: parsing response: %w", err)
	}
	if resp.Error != nil {
		return "", classifyError(resp.Error.Message)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("LLM error: empty response from provider")
	}
	return resp.Choices[0].Message.Content, nil
}
