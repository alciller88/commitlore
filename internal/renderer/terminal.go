package renderer

import (
	"fmt"
	"strings"

	"github.com/alciller88/commitlore/internal/styles"
)

var ansiCodes = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
}

const ansiReset = "\033[0m"
const ansiBold = "\033[1m"

func renderTerminal(content string, style styles.Style) string {
	tc := style.Terminal
	content = colorHeaders(content, tc)
	content = colorItems(content, tc)
	return addDecorators(content, tc)
}

func colorHeaders(s string, tc styles.Terminal) string {
	headerCode := lookupANSI(tc.Colors.Header)
	s = colorMainHeader(s, headerCode)
	s = colorSectionHeaders(s, headerCode)
	return s
}

func colorMainHeader(s, code string) string {
	if idx := strings.Index(s, "# "); idx == 0 || (idx > 0 && s[idx-1] == '\n') {
		end := strings.Index(s[idx:], "\n")
		if end < 0 {
			end = len(s) - idx
		}
		header := s[idx : idx+end]
		s = strings.Replace(s, header, ansiBold+code+header+ansiReset, 1)
	}
	return s
}

func colorSectionHeaders(s, code string) string {
	return strings.ReplaceAll(s, "\n## ", "\n"+code+ansiBold+"## ")
}

func colorItems(s string, tc styles.Terminal) string {
	return strings.ReplaceAll(s, "\n- ", ansiReset+"\n- ")
}

func addDecorators(s string, tc styles.Terminal) string {
	if tc.Decorators.Separator != "" {
		footerCode := lookupANSI(tc.Colors.Footer)
		s += fmt.Sprintf("\n%s%s%s\n", footerCode, tc.Decorators.Separator, ansiReset)
	}
	return s
}

func lookupANSI(name string) string {
	if code, ok := ansiCodes[strings.ToLower(name)]; ok {
		return code
	}
	return ""
}
