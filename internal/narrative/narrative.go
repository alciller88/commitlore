package narrative

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/styles"
)

// Generate produces narrative text from a Changelog using the given Style.
func Generate(cl changelog.Changelog, style styles.Style) (string, error) {
	data := buildTemplateData(cl)
	text, err := renderWithStyle(data, style)
	if err != nil {
		return "", err
	}
	return ApplyVocabulary(text, style.Vocabulary), nil
}

func renderWithStyle(data templateData, style styles.Style) (string, error) {
	var buf bytes.Buffer

	if err := writeHeader(&buf, style.Templates.Header); err != nil {
		return "", err
	}

	for _, g := range data.Groups {
		writeGroupHeader(&buf, g.Label)
		if err := writeGroupCommits(&buf, g, style); err != nil {
			return "", err
		}
	}

	if err := writeFooter(&buf, style.Templates.Footer); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func writeHeader(buf *bytes.Buffer, header string) error {
	return executeInline(buf, header, nil)
}

func writeFooter(buf *bytes.Buffer, footer string) error {
	if footer == "" {
		return nil
	}
	buf.WriteString("\n")
	return executeInline(buf, footer, nil)
}

func writeGroupHeader(buf *bytes.Buffer, label string) {
	fmt.Fprintf(buf, "\n## %s\n", label)
}

func writeGroupCommits(buf *bytes.Buffer, g templateGroup, style styles.Style) error {
	tmplStr := selectCommitTemplate(g.Label, style)
	for _, c := range g.Commits {
		if err := executeInline(buf, tmplStr, c); err != nil {
			return err
		}
		buf.WriteString("\n")
	}
	return nil
}

func selectCommitTemplate(label string, style styles.Style) string {
	switch label {
	case "Breaking Changes":
		return style.Templates.Breaking
	case "Bug Fixes":
		return style.Templates.Fix
	default:
		return style.Templates.Feature
	}
}

func executeInline(buf *bytes.Buffer, tmplStr string, data interface{}) error {
	funcMap := template.FuncMap{
		"short": shortHash,
		"upper": strings.ToUpper,
	}
	tmpl, err := template.New("inline").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	return tmpl.Execute(buf, data)
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}
