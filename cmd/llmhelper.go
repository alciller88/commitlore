package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/llm"
)

const llmTimeout = 30 * time.Second

// resolveLLMProvider returns the provider name from flag or env var.
func resolveLLMProvider(flagValue string) string {
	if flagValue != "" && flagValue != "none" {
		return flagValue
	}
	if env := os.Getenv("COMMITLORE_LLM_PROVIDER"); env != "" {
		return env
	}
	return "none"
}

// resolveLLMBaseURL returns the base URL from flag or env var.
func resolveLLMBaseURL(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return os.Getenv("COMMITLORE_LLM_BASE_URL")
}

// resolveLLMModel returns the model from flag or env var.
func resolveLLMModel(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return os.Getenv("COMMITLORE_LLM_MODEL")
}

// enrichWithLLM passes text through an LLM for enrichment.
// On any error (timeout, auth, network), it logs a warning and returns
// the original text. The LLM never causes a command failure.
func enrichWithLLM(provider, baseURL, model, llmPrompt, text string) string {
	if provider == "none" || llmPrompt == "" {
		return text
	}

	apiKey := os.Getenv("COMMITLORE_LLM_API_KEY")
	p, err := llm.New(provider, apiKey, baseURL, model)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %s. Using template output.\n", err)
		return text
	}

	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()

	result, err := p.Enrich(ctx, llmPrompt, text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %s. Using template output.\n", err)
		return text
	}

	return result
}
