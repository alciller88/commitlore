package changelog

import (
	"github.com/alciller88/commitlore/internal/git"
)

// groupOrder defines the display order for commit type groups.
var groupOrder = []CommitType{
	TypeBreaking,
	TypeFeat,
	TypeFix,
	TypeRefactor,
	TypeDocs,
	TypeTest,
	TypeChore,
	TypeOther,
}

// GroupCommits converts raw commits into a Changelog grouped by type.
func GroupCommits(commits []git.Commit) Changelog {
	parsed := parseAll(commits)
	grouped := buildGroupMap(parsed)
	return toChangelog(grouped)
}

func parseAll(commits []git.Commit) []ParsedCommit {
	result := make([]ParsedCommit, 0, len(commits))
	for _, c := range commits {
		result = append(result, ParsedCommit{
			Hash:    c.Hash,
			Author:  c.Author,
			Email:   c.Email,
			Date:    c.Date,
			Type:    ParseType(c.Message),
			Message: c.Message,
		})
	}
	return result
}

func buildGroupMap(commits []ParsedCommit) map[CommitType][]ParsedCommit {
	groups := make(map[CommitType][]ParsedCommit)
	for _, c := range commits {
		groups[c.Type] = append(groups[c.Type], c)
	}
	return groups
}

func toChangelog(groups map[CommitType][]ParsedCommit) Changelog {
	var cl Changelog
	for _, t := range groupOrder {
		if commits, ok := groups[t]; ok {
			cl.Groups = append(cl.Groups, ChangelogGroup{
				Type:    t,
				Commits: commits,
			})
		}
	}
	return cl
}
