package llm

import (
	"strings"
	"unicode"
)

const maxRepoDataLen = 500

// SanitizeRepoData truncates to 500 characters and escapes control characters.
// Must be called before passing any repository data to an LLM.
func SanitizeRepoData(s string) string {
	s = escapeControlChars(s)
	if len(s) > maxRepoDataLen {
		s = s[:maxRepoDataLen]
	}
	return s
}

func escapeControlChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
