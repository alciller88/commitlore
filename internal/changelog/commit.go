package changelog

import (
	"strings"
	"time"
)

// CommitType represents the category of a commit.
type CommitType string

const (
	TypeFeat     CommitType = "feat"
	TypeFix      CommitType = "fix"
	TypeChore    CommitType = "chore"
	TypeDocs     CommitType = "docs"
	TypeTest     CommitType = "test"
	TypeRefactor CommitType = "refactor"
	TypeBreaking CommitType = "breaking"
	TypeOther    CommitType = "other"
)

// ParsedCommit is a commit with its type extracted from the message.
type ParsedCommit struct {
	Hash    string
	Author  string
	Email   string
	Date    time.Time
	Type    CommitType
	Message string
}

// ChangelogGroup holds commits grouped by their type.
type ChangelogGroup struct {
	Type    CommitType
	Commits []ParsedCommit
}

// Changelog is an ordered collection of groups ready for rendering.
type Changelog struct {
	Groups []ChangelogGroup
}

// ParseType extracts the commit type from a Conventional Commits message.
// Falls back to inference from keywords when the format is not followed.
func ParseType(message string) CommitType {
	if t, ok := parseConventional(message); ok {
		return t
	}
	return inferType(message)
}

func parseConventional(message string) (CommitType, bool) {
	lower := strings.ToLower(message)

	if strings.HasPrefix(lower, "breaking change") {
		return TypeBreaking, true
	}

	idx := strings.Index(lower, ":")
	if idx < 1 {
		return "", false
	}

	prefix := strings.TrimSpace(lower[:idx])
	prefix = strings.TrimSuffix(prefix, "!")

	// Handle scoped prefixes like "feat(api)"
	if parenIdx := strings.Index(prefix, "("); parenIdx > 0 {
		prefix = prefix[:parenIdx]
	}

	return matchPrefix(prefix)
}

func matchPrefix(prefix string) (CommitType, bool) {
	switch prefix {
	case "feat", "feature":
		return TypeFeat, true
	case "fix", "bugfix":
		return TypeFix, true
	case "chore":
		return TypeChore, true
	case "docs", "doc":
		return TypeDocs, true
	case "test", "tests":
		return TypeTest, true
	case "refactor":
		return TypeRefactor, true
	default:
		return "", false
	}
}

var inferenceRules = []struct {
	keywords []string
	commitType CommitType
}{
	{[]string{"breaking", "removed", "incompatible"}, TypeBreaking},
	{[]string{"add", "new", "implement", "introduce"}, TypeFeat},
	{[]string{"fix", "bug", "patch", "resolve", "correct"}, TypeFix},
	{[]string{"refactor", "restructure", "reorganize", "clean up"}, TypeRefactor},
	{[]string{"doc", "readme", "comment"}, TypeDocs},
	{[]string{"test", "spec", "coverage"}, TypeTest},
	{[]string{"chore", "bump", "update dep", "ci", "build"}, TypeChore},
}

func inferType(message string) CommitType {
	lower := strings.ToLower(message)
	for _, rule := range inferenceRules {
		for _, kw := range rule.keywords {
			if strings.Contains(lower, kw) {
				return rule.commitType
			}
		}
	}
	return TypeOther
}
