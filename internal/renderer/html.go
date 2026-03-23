package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
	goldhtml "github.com/yuin/goldmark/renderer/html"

	"github.com/alciller88/commitlore/assets"
	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/styles"
)

func renderChangelogHTML(content string, cl changelog.Changelog, style styles.Style) (string, error) {
	if style.HTMLTemplate != "" {
		return renderCustomChangelogHTML(content, cl, style)
	}
	return renderDefaultChangelogHTML(content, cl, style)
}

func renderDefaultChangelogHTML(content string, cl changelog.Changelog, style styles.Style) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("<div class=\"narrative\">\n")
	writeNarrativeHTML(&buf, content)
	buf.WriteString("</div>\n")

	buf.WriteString("<div class=\"data-section\">\n")
	dataBody, err := buildChangelogBody(cl)
	if err != nil {
		return "", err
	}
	buf.WriteString(dataBody)
	buf.WriteString("</div>\n")

	return renderHTMLPage("Changelog — commitlore", buf.String(), style)
}

func buildChangelogBody(cl changelog.Changelog) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("<h1>Changelog</h1>\n")
	for _, g := range cl.Groups {
		if err := writeGroupHTML(&buf, g); err != nil {
			return "", err
		}
	}
	writeHTMLFooter(&buf)
	return buf.String(), nil
}

func writeGroupHTML(buf *bytes.Buffer, g changelog.ChangelogGroup) error {
	fmt.Fprintf(buf, "<h2>%s</h2>\n<ul>\n", template.HTMLEscapeString(groupTitle(g.Type)))
	for _, c := range g.Commits {
		writeCommitHTML(buf, c)
	}
	buf.WriteString("</ul>\n")
	return nil
}

func writeCommitHTML(buf *bytes.Buffer, c changelog.ParsedCommit) {
	badge := typeBadgeClass(c.Type)
	fmt.Fprintf(buf,
		"<li><span class=\"type-badge %s\">%s</span> %s "+
			"<span class=\"hash\">%s</span> — "+
			"<span class=\"author\">%s</span> "+
			"<span class=\"date\">%s</span></li>\n",
		badge, template.HTMLEscapeString(string(c.Type)),
		template.HTMLEscapeString(c.Message),
		shortHash(c.Hash),
		template.HTMLEscapeString(c.Author),
		c.Date.Format("2006-01-02"),
	)
}

func renderStoryHTML(content string, ch git.Chronology, style styles.Style) (string, error) {
	if style.HTMLTemplate != "" {
		return renderCustomStoryHTML(content, ch, style)
	}
	return renderDefaultStoryHTML(content, ch, style)
}

func renderDefaultStoryHTML(content string, ch git.Chronology, style styles.Style) (string, error) {
	var buf bytes.Buffer

	buf.WriteString("<div class=\"narrative\">\n")
	writeNarrativeHTML(&buf, content)
	buf.WriteString("</div>\n")

	buf.WriteString("<div class=\"data-section\">\n")
	buf.WriteString(buildStoryBody(ch))
	buf.WriteString("</div>\n")

	return renderHTMLPage("Story — commitlore", buf.String(), style)
}

func buildStoryBody(ch git.Chronology) string {
	var buf bytes.Buffer
	buf.WriteString("<h1>Repository Story</h1>\n")
	writeStoryIntroHTML(&buf, ch)
	writeStoryTagsHTML(&buf, ch)
	writeStoryPeaksHTML(&buf, ch)
	writeStoryContributorsHTML(&buf, ch)
	writeHTMLFooter(&buf)
	return buf.String()
}

func writeStoryIntroHTML(buf *bytes.Buffer, ch git.Chronology) {
	if ch.TotalCommits == 0 {
		return
	}
	fmt.Fprintf(buf,
		"<p>Started on <strong>%s</strong> by <span class=\"author\">%s</span>. "+
			"Total commits: <strong>%d</strong>. Contributors: <strong>%d</strong>.</p>\n",
		ch.FirstCommit.Date.Format("2006-01-02"),
		template.HTMLEscapeString(ch.FirstCommit.Author),
		ch.TotalCommits, len(ch.Contributors),
	)
}

func writeStoryTagsHTML(buf *bytes.Buffer, ch git.Chronology) {
	if len(ch.Tags) == 0 {
		return
	}
	buf.WriteString("<h2>Milestones</h2>\n<ul>\n")
	for _, t := range ch.Tags {
		fmt.Fprintf(buf,
			"<li><strong>%s</strong> — <span class=\"date\">%s</span> <span class=\"hash\">%s</span></li>\n",
			template.HTMLEscapeString(t.Name), t.Date.Format("2006-01-02"), shortHash(t.Hash))
	}
	buf.WriteString("</ul>\n")
}

func writeStoryPeaksHTML(buf *bytes.Buffer, ch git.Chronology) {
	if len(ch.Peaks) == 0 {
		return
	}
	maxCount := ch.Peaks[0].Count
	buf.WriteString("<h2>Activity Peaks</h2>\n<ul>\n")
	for _, p := range ch.Peaks {
		barWidth := peakBarWidth(p.Count, maxCount)
		fmt.Fprintf(buf,
			"<li><span class=\"peak-bar\" style=\"width:%dpx\"></span> %s — <strong>%d</strong> commits</li>\n",
			barWidth, template.HTMLEscapeString(p.Month), p.Count)
	}
	buf.WriteString("</ul>\n")
}

func writeStoryContributorsHTML(buf *bytes.Buffer, ch git.Chronology) {
	if len(ch.Contributors) == 0 {
		return
	}
	buf.WriteString("<h2>Contributors</h2>\n<ul>\n")
	for _, c := range ch.Contributors {
		fmt.Fprintf(buf,
			"<li><span class=\"author\">%s</span> — joined <span class=\"date\">%s</span></li>\n",
			template.HTMLEscapeString(c.Name), c.Date.Format("2006-01-02"))
	}
	buf.WriteString("</ul>\n")
}

func renderHTMLPage(title, body string, style styles.Style) (string, error) {
	css := buildCSS(style.Theme)
	html := buildHTMLDocument(title, css, body, style.Theme)
	return html, nil
}

func buildHTMLDocument(title, css, body string, theme styles.Theme) string {
	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	fmt.Fprintf(&buf, "<meta charset=\"UTF-8\">\n<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	fmt.Fprintf(&buf, "<title>%s</title>\n", template.HTMLEscapeString(title))
	fmt.Fprintf(&buf, "<style>\n%s\n</style>\n", css)
	buf.WriteString("</head>\n<body>\n")
	writeSiteHeader(&buf, theme)
	buf.WriteString("<div class=\"content\">\n")
	buf.WriteString(body)
	buf.WriteString("</div>\n")
	buf.WriteString("</body>\n</html>")
	return buf.String()
}

func writeSiteHeader(buf *bytes.Buffer, theme styles.Theme) {
	if theme.HeaderImage != "" {
		fmt.Fprintf(buf, "<img class=\"header-image\" src=\"%s\" alt=\"header\">\n",
			template.HTMLEscapeString(theme.HeaderImage))
	}
	buf.WriteString("<header class=\"site-header\">\n")
	writeSiteHeaderLogo(buf, theme)
	buf.WriteString("</header>\n")
}

func writeSiteHeaderLogo(buf *bytes.Buffer, theme styles.Theme) {
	if theme.Logo != "" {
		fmt.Fprintf(buf, "<img class=\"logo\" src=\"%s\" alt=\"logo\">\n",
			template.HTMLEscapeString(theme.Logo))
		return
	}
	logoSVG := strings.Replace(assets.LogoSVG, "<svg ", `<svg width="100" height="100" `, 1)
	fmt.Fprintf(buf, "<div class=\"logo-svg\">%s</div>\n", logoSVG)
}

func markdownToHTML(s string) string {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			goldhtml.WithHardWraps(),
			goldhtml.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(s), &buf); err != nil {
		return template.HTMLEscapeString(s)
	}
	return buf.String()
}

func writeNarrativeHTML(buf *bytes.Buffer, content string) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return
	}
	buf.WriteString(markdownToHTML(trimmed))
}

func shortHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func groupTitle(t changelog.CommitType) string {
	titles := map[changelog.CommitType]string{
		changelog.TypeBreaking: "Breaking Changes",
		changelog.TypeFeat:     "Features",
		changelog.TypeFix:      "Bug Fixes",
		changelog.TypeRefactor: "Refactoring",
		changelog.TypeDocs:     "Documentation",
		changelog.TypeTest:     "Tests",
		changelog.TypeChore:    "Chores",
		changelog.TypeOther:    "Other",
	}
	if title, ok := titles[t]; ok {
		return title
	}
	return string(t)
}

func typeBadgeClass(t changelog.CommitType) string {
	switch t {
	case changelog.TypeFeat:
		return "type-feat"
	case changelog.TypeFix:
		return "type-fix"
	case changelog.TypeBreaking:
		return "type-breaking"
	default:
		return "type-other"
	}
}

func peakBarWidth(count, max int) int {
	if max == 0 {
		return 0
	}
	return (count * 200) / max
}

func writeHTMLFooter(buf *bytes.Buffer) {
	buf.WriteString("<div class=\"footer\"></div>\n")
}
