package narrative

import (
	"strings"
	"unicode"
)

// ApplyVocabulary replaces whole-word keywords in text using the vocabulary map.
// Replacements are case-insensitive and respect word boundaries.
func ApplyVocabulary(text string, vocab map[string]string) string {
	if len(vocab) == 0 {
		return text
	}
	for word, replacement := range vocab {
		text = replaceWholeWord(text, word, replacement)
	}
	return text
}

func replaceWholeWord(text, old, replacement string) string {
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
		absIdx := i + idx
		if !isWordBoundary(text, absIdx, len(old)) {
			result.WriteString(text[i : absIdx+len(old)])
			i = absIdx + len(old)
			continue
		}
		result.WriteString(text[i:absIdx])
		result.WriteString(replacement)
		i = absIdx + len(old)
	}
	return result.String()
}

func isWordBoundary(text string, start, length int) bool {
	if start > 0 && isWordChar(rune(text[start-1])) {
		return false
	}
	end := start + length
	if end < len(text) && isWordChar(rune(text[end])) {
		return false
	}
	return true
}

func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
