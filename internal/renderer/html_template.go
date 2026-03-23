package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/styles"
)

// HTMLTemplateContext provides all data to custom HTML templates.
// It is used for both changelog and story renders; fields not relevant
// to a given mode are left at their zero values so {{if .Items}} and
// {{if .Peaks}} conditionals work correctly for both contexts.
type HTMLTemplateContext struct {
	Title     string
	Content   template.HTML
	Items     []HTMLItem
	Theme     styles.Theme
	Icons     styles.Icons
	UILabels  styles.UILabels
	RepoName  string
	Generated string
	Version   string

	// Changelog temporal data
	CommitsByWeek []WeekActivity

	// Story-only fields (zero when rendering a changelog)
	TotalCommits int
	FirstAuthor  string
	FirstDate    string
	Tags         []StoryTag
	Peaks        []StoryPeak
	Contributors []StoryContributor
}

// StoryHTMLContext is an alias for HTMLTemplateContext kept for compatibility.
type StoryHTMLContext = HTMLTemplateContext

// HTMLItem represents a single commit entry for the HTML template.
type HTMLItem struct {
	Type    string
	Message string
	Hash    string
	Author  string
	Date    string
	Icon    string
}

// StoryTag represents a tag entry for story templates.
type StoryTag struct {
	Name string
	Hash string
	Date string
}

// StoryPeak represents an activity peak for story templates.
type StoryPeak struct {
	Month string
	Count int
}

// StoryContributor represents a contributor entry for story templates.
type StoryContributor struct {
	Name  string
	Email string
	Date  string
	Count int
}

// WeekActivity represents commit count for a calendar week.
type WeekActivity struct {
	Label string
	Count int
}

func groupItemsByWeek(items []HTMLItem) []WeekActivity {
	counts := make(map[string]int)
	var keys []string
	for _, item := range items {
		week := weekLabel(item.Date)
		if _, exists := counts[week]; !exists {
			keys = append(keys, week)
		}
		counts[week]++
	}
	result := make([]WeekActivity, 0, len(keys))
	for _, k := range keys {
		result = append(result, WeekActivity{Label: k, Count: counts[k]})
	}
	return result
}

func weekLabel(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

// SemverFromString returns s if it matches a semver tag (vMAJOR.MINOR.PATCH),
// otherwise returns "".
func SemverFromString(s string) string {
	s = strings.TrimSpace(s)
	matched, _ := regexp.MatchString(`^v?\d+\.\d+\.\d+`, s)
	if matched {
		return s
	}
	return ""
}

func renderCustomChangelogHTML(content string, cl changelog.Changelog, style styles.Style, repoName, version string) (string, error) {
	ctx := buildChangelogContext(content, cl, style, repoName, version)
	return executeHTMLTemplate(style.HTMLTemplateChangelog, ctx)
}

func renderCustomStoryHTML(content string, ch git.Chronology, style styles.Style, repoName string) (string, error) {
	ctx := buildStoryContext(content, ch, style, repoName)
	return executeHTMLTemplate(style.HTMLTemplateStory, ctx)
}

func buildChangelogContext(content string, cl changelog.Changelog, style styles.Style, repoName, version string) HTMLTemplateContext {
	if repoName == "" {
		repoName = "Repository"
	}
	items := extractItems(cl, style.Icons)
	return HTMLTemplateContext{
		Title:         "Changelog",
		Content:       template.HTML(markdownToHTML(strings.TrimSpace(content))),
		Items:         items,
		Theme:         withDefaults(style.Theme),
		Icons:         style.Icons,
		UILabels:      style.UILabels,
		RepoName:      repoName,
		Generated:     time.Now().Format("2 Jan 2006"),
		Version:       version,
		CommitsByWeek: groupItemsByWeek(items),
	}
}

func buildStoryContext(content string, ch git.Chronology, style styles.Style, repoName string) HTMLTemplateContext {
	if repoName == "" {
		repoName = "Repository"
	}
	tags := make([]StoryTag, 0, len(ch.Tags))
	for _, t := range ch.Tags {
		tags = append(tags, StoryTag{Name: t.Name, Hash: shortHash(t.Hash), Date: t.Date.Format("2006-01-02")})
	}
	peaks := make([]StoryPeak, 0, len(ch.Peaks))
	for _, p := range ch.Peaks {
		peaks = append(peaks, StoryPeak{Month: p.Month, Count: p.Count})
	}
	contribs := make([]StoryContributor, 0, len(ch.Contributors))
	for _, c := range ch.Contributors {
		contribs = append(contribs, StoryContributor{Name: c.Name, Email: c.Email, Date: c.Date.Format("2006-01-02"), Count: c.Count})
	}
	firstAuthor := ""
	firstDate := ""
	if ch.TotalCommits > 0 {
		firstAuthor = ch.FirstCommit.Author
		firstDate = ch.FirstCommit.Date.Format("2006-01-02")
	}
	if firstAuthor == "" {
		firstAuthor = "an unknown contributor"
	}
	return HTMLTemplateContext{
		Title:        "Repository Story",
		Content:      template.HTML(markdownToHTML(strings.TrimSpace(content))),
		Theme:        withDefaults(style.Theme),
		Icons:        style.Icons,
		UILabels:     style.UILabels,
		RepoName:     repoName,
		Generated:    time.Now().Format("2 Jan 2006"),
		Version:      "",
		TotalCommits: ch.TotalCommits,
		FirstAuthor:  firstAuthor,
		FirstDate:    firstDate,
		Tags:         tags,
		Peaks:        peaks,
		Contributors: contribs,
	}
}

func extractItems(cl changelog.Changelog, icons styles.Icons) []HTMLItem {
	var items []HTMLItem
	for _, g := range cl.Groups {
		for _, c := range g.Commits {
			items = append(items, HTMLItem{
				Type:    string(c.Type),
				Message: c.Message,
				Hash:    shortHash(c.Hash),
				Author:  c.Author,
				Date:    c.Date.Format("2006-01-02"),
				Icon:    iconForType(c.Type, icons),
			})
		}
	}
	return items
}

func iconForType(t changelog.CommitType, icons styles.Icons) string {
	switch t {
	case changelog.TypeFeat:
		return icons.Feature
	case changelog.TypeFix:
		return icons.Fix
	case changelog.TypeBreaking:
		return icons.Breaking
	case changelog.TypeChore:
		return icons.Chore
	case changelog.TypeDocs:
		return icons.Docs
	case changelog.TypeTest:
		return icons.Test
	default:
		return icons.Bullet
	}
}

func executeHTMLTemplate(tmplStr string, data interface{}) (string, error) {
	funcMap := template.FuncMap{
		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"initials": templateInitials,
		"add":      func(a, b int) int { return a + b },
		"mul":      func(a, b float64) float64 { return a * b },
		"divf": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"divi": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"safeCSS":  func(s string) template.CSS { return template.CSS(s) },
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"safeAttr": func(s string) template.HTMLAttr { return template.HTMLAttr(s) },
	}
	tmpl, err := template.New("html").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parsing html_template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing html_template: %w", err)
	}
	return buf.String(), nil
}

func templateInitials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}
	result := ""
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(string([]rune(p)[0]))
		}
		if len(result) >= 2 {
			break
		}
	}
	return result
}
