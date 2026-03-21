package narrative

import (
	"strings"
)

// ApplyVocabulary replaces keywords in text using the vocabulary map.
// Replacements are case-insensitive but preserve surrounding text.
func ApplyVocabulary(text string, vocab map[string]string) string {
	if len(vocab) == 0 {
		return text
	}
	for word, replacement := range vocab {
		text = replaceInsensitive(text, word, replacement)
	}
	return text
}

func replaceInsensitive(text, old, replacement string) string {
	lower := strings.ToLower(text)
	oldLower := strings.ToLower(old)
	var result strings.Builder
	i := 0
	for {
		idx := strings.Index(lower[i:], oldLower)
		if idx < 0 {
			result.WriteString(text[i:])
			break
		}
		result.WriteString(text[i : i+idx])
		result.WriteString(replacement)
		i += idx + len(old)
	}
	return result.String()
}
