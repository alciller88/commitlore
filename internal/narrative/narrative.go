package narrative

import (
	"encoding/json"
	"fmt"

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

// RenderOptions configures how a changelog is rendered.
type RenderOptions struct {
	Style  string
	Format Format
}

// Render produces a string from a Changelog using the given options.
func Render(cl changelog.Changelog, opts RenderOptions) (string, error) {
	switch opts.Format {
	case FormatJSON:
		return renderJSON(cl)
	case FormatHTML, FormatPDF:
		return "", fmt.Errorf("format %q not implemented yet", opts.Format)
	default:
		return renderTemplate(cl, opts)
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
