package llm

import "context"

// Provider is the interface for LLM adapters.
// All providers implement a single Enrich method that takes instructions
// and data, and returns enriched text.
type Provider interface {
	Enrich(ctx context.Context, prompt, data string) (string, error)
}
