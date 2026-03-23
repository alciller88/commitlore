package renderer

import (
	"encoding/json"
	"fmt"
	"strings"

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

// HTMLTheme allows overriding the style's theme for HTML output.
type HTMLTheme struct {
	Background string
	Surface    string
	Text       string
	Primary    string
	Secondary  string
	Accent     string
	Border     string
	FontFamily string
	Mode       string
}

// Render formats the given content according to the specified format.
// Optional variadic args: opts[0] = repoName, opts[1] = version.
func Render(content string, cl changelog.Changelog, style styles.Style, format Format, opts ...string) (string, error) {
	rn, version := extractOpts(opts)
	switch format {
	case FormatJSON:
		return renderJSON(cl)
	case FormatHTML:
		return renderChangelogHTML(content, cl, style, rn, version)
	case FormatPDF:
		return "", fmt.Errorf("PDF format has been removed. Use --format html instead.")
	case FormatTerminal:
		return renderTerminal(content, style), nil
	default:
		return content, nil
	}
}

// RenderWithTheme renders with an optional theme override for HTML output.
// If override is nil, uses the style's own theme (same as Render).
// Optional variadic args: opts[0] = repoName, opts[1] = version.
func RenderWithTheme(content string, cl changelog.Changelog, style styles.Style, format Format, override *HTMLTheme, opts ...string) (string, error) {
	if override != nil {
		style = applyHTMLThemeOverride(style, override)
	}
	return Render(content, cl, style, format, opts...)
}

func extractOpts(opts []string) (repoName, version string) {
	repoName = "Repository"
	if len(opts) > 0 && opts[0] != "" {
		repoName = opts[0]
	}
	if len(opts) > 1 {
		version = opts[1]
	}
	return
}

func extractRepoName(names []string) string {
	if len(names) > 0 && names[0] != "" {
		return names[0]
	}
	return "Repository"
}

// RepoNameFromPath extracts a display name from a repo path or GitHub URL.
// "C:\Users\alcil\MyProjects\commitlore" → "commitlore"
// "owner/repo" → "repo"
// "." → "Repository"
func RepoNameFromPath(repoPath string) string {
	if repoPath == "" || repoPath == "." {
		return "Repository"
	}
	// GitHub format: owner/repo
	if strings.Contains(repoPath, "/") && !strings.Contains(repoPath, "\\") && !strings.HasPrefix(repoPath, "/") {
		parts := strings.Split(repoPath, "/")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
	}
	// Local path: extract last segment
	clean := strings.TrimRight(repoPath, "/\\")
	for i := len(clean) - 1; i >= 0; i-- {
		if clean[i] == '/' || clean[i] == '\\' {
			return clean[i+1:]
		}
	}
	return clean
}

func applyHTMLThemeOverride(style styles.Style, o *HTMLTheme) styles.Style {
	s := style
	overrideField(&s.Theme.Colors.Background, o.Background)
	overrideField(&s.Theme.Colors.Surface, o.Surface)
	overrideField(&s.Theme.Colors.Text, o.Text)
	overrideField(&s.Theme.Colors.Primary, o.Primary)
	overrideField(&s.Theme.Colors.Secondary, o.Secondary)
	overrideField(&s.Theme.Colors.Accent, o.Accent)
	overrideField(&s.Theme.Colors.Border, o.Border)
	overrideField(&s.Theme.Typography.FontFamily, o.FontFamily)
	overrideField(&s.Theme.Mode, o.Mode)
	return s
}

func overrideField(field *string, value string) {
	if value != "" {
		*field = value
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
