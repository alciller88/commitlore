package renderer

import (
	"encoding/json"
	"fmt"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/styles"
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
func Render(content string, cl changelog.Changelog, style styles.Style, format Format) (string, error) {
	switch format {
	case FormatJSON:
		return renderJSON(cl)
	case FormatHTML:
		return renderChangelogHTML(content, cl, style)
	case FormatPDF:
		return "", fmt.Errorf("PDF format has been removed. Use --format html instead.")
	case FormatTerminal:
		return renderTerminal(content, style), nil
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
