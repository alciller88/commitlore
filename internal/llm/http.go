package llm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type responseParser func([]byte) (string, error)

func executeAndParse(client *http.Client, req *http.Request, parse responseParser) (string, error) {
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM error: request failed: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("LLM error: reading response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("LLM error: authentication failed — check COMMITLORE_LLM_API_KEY")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("LLM error: provider returned status %d", resp.StatusCode)
	}

	return parse(body)
}

func classifyError(msg string) error {
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "auth") || strings.Contains(lower, "api key") || strings.Contains(lower, "invalid") {
		return fmt.Errorf("LLM error: authentication failed — check COMMITLORE_LLM_API_KEY")
	}
	return fmt.Errorf("LLM error: %s", msg)
}
