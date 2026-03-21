package renderer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alciller88/commitlore/internal/changelog"
)

// Format represents an output format.
type Format string

const (
	FormatTerminal Format = "terminal"
	FormatMD       Format = "md"
	FormatJSON     Format = "json"
	FormatHTML     Format = "html"
	FormatPDF      Format = "pdf"
)

// Render formats the given content according to the specified format.
// For terminal and md, content is the narrative text.
// For json, the changelog is serialized directly.
func Render(content string, cl changelog.Changelog, format Format) (string, error) {
	switch format {
	case FormatJSON:
		return renderJSON(cl)
	case FormatHTML:
		return renderChangelogHTML(content, cl)
	case FormatPDF:
		return "", fmt.Errorf("PDF format has been removed. Use --format html instead.")
	case FormatTerminal:
		return addANSIColors(content), nil
	default:
		return content, nil
	}
}

func renderJSON(cl changelog.Changelog) (string, error) {
	data := toJSONData(cl)
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding JSON: %w", err)
	}
	return string(out), nil
}

func addANSIColors(s string) string {
	s = strings.ReplaceAll(s, "# Changelog", "\033[1m# Changelog\033[0m")
	s = strings.ReplaceAll(s, "## ", "\033[1;36m## ")
	s = strings.ReplaceAll(s, "\n- ", "\033[0m\n- ")
	return s
}
