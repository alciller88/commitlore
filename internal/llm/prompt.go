package llm

// BuildPrompt wraps repo data between delimiters as required by SPEC §15.
// The llmPrompt contains instructions; data contains sanitized repo content.
func BuildPrompt(llmPrompt, data string) string {
	return llmPrompt + "\n\n---DATA START---\n" + data + "\n---DATA END---"
}
