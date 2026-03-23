package renderer

import (
	"encoding/json"
	"fmt"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/styles"
)

// RenderStoryWithTheme renders story with an optional theme override for HTML output.
func RenderStoryWithTheme(content string, ch git.Chronology, style styles.Style, format Format, override *HTMLTheme, repoName ...string) (string, error) {
	if override != nil {
		style = applyHTMLThemeOverride(style, override)
	}
	return RenderStory(content, ch, style, format, repoName...)
}

// RenderStory formats story content according to the specified format.
func RenderStory(content string, ch git.Chronology, style styles.Style, format Format, repoName ...string) (string, error) {
	rn := extractRepoName(repoName)
	switch format {
	case FormatJSON:
		return renderStoryJSON(ch)
	case FormatHTML:
		return renderStoryHTML(content, ch, style, rn)
	case FormatPDF:
		return "", fmt.Errorf("PDF format has been removed. Use --format html instead.")
	case FormatTerminal:
		return renderTerminal(content, style), nil
	default:
		return content, nil
	}
}

func renderStoryJSON(ch git.Chronology) (string, error) {
	data := toStoryJSON(ch)
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding story JSON: %w", err)
	}
	return string(out), nil
}

type jsonStory struct {
	FirstCommit  jsonStoryCommit   `json:"first_commit"`
	TotalCommits int               `json:"total_commits"`
	Tags         []jsonTag         `json:"tags"`
	Peaks        []jsonPeak        `json:"activity_peaks"`
	Contributors []jsonContributor `json:"contributors"`
}

type jsonStoryCommit struct {
	Hash   string `json:"hash"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

type jsonTag struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Date string `json:"date"`
}

type jsonPeak struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type jsonContributor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

func toStoryJSON(ch git.Chronology) jsonStory {
	return jsonStory{
		FirstCommit:  toJSONStoryCommit(ch.FirstCommit),
		TotalCommits: ch.TotalCommits,
		Tags:         toJSONTags(ch.Tags),
		Peaks:        toJSONPeaks(ch.Peaks),
		Contributors: toJSONContributors(ch.Contributors),
	}
}

func toJSONStoryCommit(c git.Commit) jsonStoryCommit {
	return jsonStoryCommit{
		Hash:   c.Hash,
		Author: c.Author,
		Date:   c.Date.Format("2006-01-02"),
	}
}

func toJSONTags(tags []git.Tag) []jsonTag {
	result := make([]jsonTag, 0, len(tags))
	for _, t := range tags {
		result = append(result, jsonTag{
			Name: t.Name, Hash: t.Hash, Date: t.Date.Format("2006-01-02"),
		})
	}
	return result
}

func toJSONPeaks(peaks []git.ActivityPeak) []jsonPeak {
	result := make([]jsonPeak, 0, len(peaks))
	for _, p := range peaks {
		result = append(result, jsonPeak{Month: p.Month, Count: p.Count})
	}
	return result
}

func toJSONContributors(entries []git.ContributorEntry) []jsonContributor {
	result := make([]jsonContributor, 0, len(entries))
	for _, c := range entries {
		result = append(result, jsonContributor{
			Name: c.Name, Email: c.Email, Date: c.Date.Format("2006-01-02"),
		})
	}
	return result
}
