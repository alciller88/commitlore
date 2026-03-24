package narrative

import (
	"github.com/alciller88/commitlore/internal/changelog"
)

// templateData is the struct passed to style templates.
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

func groupLabel(t changelog.CommitType) string {
	if label, ok := typeLabels[t]; ok {
		return label
	}
	return string(t)
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
