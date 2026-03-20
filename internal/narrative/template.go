package narrative

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/alciller88/commitlore/internal/changelog"
)

//go:embed templates/*.tmpl
var embeddedTemplates embed.FS

// templateData is the struct passed to Go templates.
type templateData struct {
	Groups []templateGroup
}

type templateGroup struct {
	Label   string
	Commits []templateCommit
}

type templateCommit struct {
	Hash    string
	Author  string
	Email   string
	Date    string
	Type    string
	Message string
}

var validStyles = map[string]bool{
	"formal":     true,
	"patchnotes": true,
	"ironic":     true,
	"epic":       true,
}

func renderTemplate(cl changelog.Changelog, opts RenderOptions) (string, error) {
	style := resolveStyle(opts.Style)
	tmplStr, err := loadTemplate(style)
	if err != nil {
		return "", err
	}
	return executeTemplate(tmplStr, cl, opts.Format)
}

func resolveStyle(style string) string {
	if style == "" || !validStyles[style] {
		return "formal"
	}
	return style
}

func loadTemplate(style string) (string, error) {
	path := "templates/" + style + ".tmpl"
	data, err := embeddedTemplates.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("loading template %s: %w", style, err)
	}
	return string(data), nil
}

func executeTemplate(tmplStr string, cl changelog.Changelog, format Format) (string, error) {
	funcMap := template.FuncMap{
		"short": shortHash,
		"date":  func(s string) string { return s },
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New("changelog").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	data := buildTemplateData(cl)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	output := buf.String()
	if format == FormatTerminal {
		output = addANSIColors(output)
	}
	return output, nil
}

func buildTemplateData(cl changelog.Changelog) templateData {
	var groups []templateGroup
	for _, g := range cl.Groups {
		groups = append(groups, templateGroup{
			Label:   groupLabel(g.Type),
			Commits: toTemplateCommits(g.Commits),
		})
	}
	return templateData{Groups: groups}
}

func toTemplateCommits(commits []changelog.ParsedCommit) []templateCommit {
	result := make([]templateCommit, 0, len(commits))
	for _, c := range commits {
		result = append(result, templateCommit{
			Hash:    c.Hash,
			Author:  c.Author,
			Email:   c.Email,
			Date:    c.Date.Format("2006-01-02"),
			Type:    string(c.Type),
			Message: c.Message,
		})
	}
	return result
}

var typeLabels = map[changelog.CommitType]string{
	changelog.TypeBreaking: "Breaking Changes",
	changelog.TypeFeat:     "Features",
	changelog.TypeFix:      "Bug Fixes",
	changelog.TypeRefactor: "Refactoring",
	changelog.TypeDocs:     "Documentation",
	changelog.TypeTest:     "Tests",
	changelog.TypeChore:    "Chores",
	changelog.TypeOther:    "Other",
}

func groupLabel(t changelog.CommitType) string {
	if label, ok := typeLabels[t]; ok {
		return label
	}
	return string(t)
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func addANSIColors(s string) string {
	s = strings.ReplaceAll(s, "# Changelog", "\033[1m# Changelog\033[0m")
	s = strings.ReplaceAll(s, "## ", "\033[1;36m## ")
	s = strings.ReplaceAll(s, "\n- ", "\033[0m\n- ")
	return s
}
