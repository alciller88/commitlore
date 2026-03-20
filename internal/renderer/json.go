package renderer

import (
	"github.com/alciller88/commitlore/internal/changelog"
)

type jsonChangelog struct {
	Groups []jsonGroup `json:"groups"`
}

type jsonGroup struct {
	Type    string       `json:"type"`
	Commits []jsonCommit `json:"commits"`
}

type jsonCommit struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Email   string `json:"email"`
	Date    string `json:"date"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func toJSONData(cl changelog.Changelog) jsonChangelog {
	var groups []jsonGroup
	for _, g := range cl.Groups {
		groups = append(groups, jsonGroup{
			Type:    string(g.Type),
			Commits: toJSONCommits(g.Commits),
		})
	}
	return jsonChangelog{Groups: groups}
}

func toJSONCommits(commits []changelog.ParsedCommit) []jsonCommit {
	result := make([]jsonCommit, 0, len(commits))
	for _, c := range commits {
		result = append(result, jsonCommit{
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
