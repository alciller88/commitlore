package renderer

import (
	"fmt"
	"regexp"
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

var detailPattern = regexp.MustCompile(`\s*[\(\[][^)\]]*[\)\]]`)

func renderTerminal(content string, style styles.Style) string {
	tc := style.Terminal
	content = applyBulletAndIndent(content, tc)
	content = colorSections(content, tc)
	content = applyDensity(content, tc)
	return addDecorators(content, tc)
}

func applyBulletAndIndent(s string, tc styles.Terminal) string {
	bullet := effectiveBullet(tc)
	indent := effectiveIndent(tc)
	s = strings.ReplaceAll(s, "\n- ", "\n"+bullet+" ")
	s = strings.ReplaceAll(s, "\n  - ", "\n"+indent+bullet+" ")
	s = strings.ReplaceAll(s, "\n  ✨", "\n"+indent+"✨")
	s = strings.ReplaceAll(s, "\n  🛡", "\n"+indent+"🛡")
	s = strings.ReplaceAll(s, "\n  💀", "\n"+indent+"💀")
	s = strings.ReplaceAll(s, "\n  ⚜", "\n"+indent+"⚜")
	s = strings.ReplaceAll(s, "\n  ⚔", "\n"+indent+"⚔")
	return s
}

func effectiveBullet(tc styles.Terminal) string {
	if tc.Decorators.Bullet != "" {
		return tc.Decorators.Bullet
	}
	return "-"
}

func effectiveIndent(tc styles.Terminal) string {
	if tc.Decorators.Indent != "" {
		return tc.Decorators.Indent
	}
	return "  "
}

func colorSections(s string, tc styles.Terminal) string {
	sections := strings.Split(s, "\n## ")
	if len(sections) <= 1 {
		return colorMainHeader(s, lookupANSI(tc.Colors.Header))
	}
	result := colorMainHeader(sections[0], lookupANSI(tc.Colors.Header))
	for _, section := range sections[1:] {
		result += colorOneSection(section, tc)
	}
	return result
}

func colorOneSection(section string, tc styles.Terminal) string {
	headerCode := lookupANSI(tc.Colors.Header)
	itemColor := sectionItemColor(section, tc)
	lines := strings.SplitN(section, "\n", 2)
	header := "\n" + headerCode + ansiBold + "## " + lines[0] + ansiReset
	if len(lines) < 2 {
		return header
	}
	return header + "\n" + colorLines(lines[1], itemColor)
}

func sectionItemColor(section string, tc styles.Terminal) string {
	lower := strings.ToLower(section)
	if strings.HasPrefix(lower, "bug") || strings.HasPrefix(lower, "fix") {
		return lookupANSI(tc.Colors.Fix)
	}
	if strings.HasPrefix(lower, "break") {
		return lookupANSI(tc.Colors.Breaking)
	}
	return lookupANSI(tc.Colors.Feature)
}

func colorLines(body, color string) string {
	if color == "" {
		return body
	}
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = color + line + ansiReset
		}
	}
	return strings.Join(lines, "\n")
}

func colorMainHeader(s, code string) string {
	idx := strings.Index(s, "# ")
	if idx < 0 {
		return s
	}
	if idx > 0 && s[idx-1] != '\n' {
		return s
	}
	end := strings.Index(s[idx:], "\n")
	if end < 0 {
		end = len(s) - idx
	}
	header := s[idx : idx+end]
	return strings.Replace(s, header, ansiBold+code+header+ansiReset, 1)
}

func applyDensity(s string, tc styles.Terminal) string {
	if tc.Density != "compact" {
		return s
	}
	return detailPattern.ReplaceAllString(s, "")
}

func addDecorators(s string, tc styles.Terminal) string {
	if tc.Decorators.Separator == "" {
		return s
	}
	footerCode := lookupANSI(tc.Colors.Footer)
	return s + fmt.Sprintf("\n%s%s%s\n", footerCode, tc.Decorators.Separator, ansiReset)
}

func lookupANSI(name string) string {
	if code, ok := ansiCodes[strings.ToLower(name)]; ok {
		return code
	}
	return ""
}
