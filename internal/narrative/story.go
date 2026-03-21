package narrative

import (
	"bytes"
	"fmt"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/styles"
)

// storyData is the top-level template context for story_intro and story_footer.
type storyData struct {
	FirstCommitDate   string
	FirstAuthor       string
	TotalCommits      int
	TotalContributors int
}

// milestoneData is the template context for story_milestone.
type milestoneData struct {
	Name string
	Hash string
	Date string
}

// peakData is the template context for story_peak.
type peakData struct {
	Month string
	Count int
}

// contributorData is the template context for story_contributor.
type contributorData struct {
	Name  string
	Email string
	Date  string
}

// GenerateStory produces a narrative from a repository chronology using the given style.
func GenerateStory(ch git.Chronology, style styles.Style) (string, error) {
	var buf bytes.Buffer

	if err := writeStoryIntro(&buf, ch, style); err != nil {
		return "", err
	}
	if err := writeStoryMilestones(&buf, ch, style); err != nil {
		return "", err
	}
	if err := writeStoryPeaks(&buf, ch, style); err != nil {
		return "", err
	}
	if err := writeStoryContributors(&buf, ch, style); err != nil {
		return "", err
	}
	if err := writeStoryFooter(&buf, ch, style); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func buildStoryData(ch git.Chronology) storyData {
	return storyData{
		FirstCommitDate:   ch.FirstCommit.Date.Format("2006-01-02"),
		FirstAuthor:       ch.FirstCommit.Author,
		TotalCommits:      ch.TotalCommits,
		TotalContributors: len(ch.Contributors),
	}
}

func writeStoryIntro(buf *bytes.Buffer, ch git.Chronology, style styles.Style) error {
	if style.Templates.StoryIntro == "" {
		return nil
	}
	return executeInline(buf, style.Templates.StoryIntro, buildStoryData(ch))
}

func writeStoryMilestones(buf *bytes.Buffer, ch git.Chronology, style styles.Style) error {
	if len(ch.Tags) == 0 || style.Templates.StoryMilestone == "" {
		return nil
	}
	fmt.Fprintf(buf, "\n\n## Milestones\n")
	for _, tag := range ch.Tags {
		if err := renderMilestone(buf, tag, style); err != nil {
			return err
		}
	}
	return nil
}

func renderMilestone(buf *bytes.Buffer, tag git.Tag, style styles.Style) error {
	data := milestoneData{
		Name: tag.Name,
		Hash: tag.Hash,
		Date: tag.Date.Format("2006-01-02"),
	}
	buf.WriteString("\n")
	return executeInline(buf, style.Templates.StoryMilestone, data)
}

func writeStoryPeaks(buf *bytes.Buffer, ch git.Chronology, style styles.Style) error {
	if len(ch.Peaks) == 0 || style.Templates.StoryPeak == "" {
		return nil
	}
	fmt.Fprintf(buf, "\n\n## Activity Peaks\n")
	for _, peak := range ch.Peaks {
		buf.WriteString("\n")
		if err := executeInline(buf, style.Templates.StoryPeak, peakData(peak)); err != nil {
			return err
		}
	}
	return nil
}

func writeStoryContributors(buf *bytes.Buffer, ch git.Chronology, style styles.Style) error {
	if len(ch.Contributors) == 0 || style.Templates.StoryContributor == "" {
		return nil
	}
	fmt.Fprintf(buf, "\n\n## Contributors\n")
	for _, c := range ch.Contributors {
		if err := renderContributor(buf, c, style); err != nil {
			return err
		}
	}
	return nil
}

func renderContributor(buf *bytes.Buffer, c git.ContributorEntry, style styles.Style) error {
	data := contributorData{
		Name:  c.Name,
		Email: c.Email,
		Date:  c.Date.Format("2006-01-02"),
	}
	buf.WriteString("\n")
	return executeInline(buf, style.Templates.StoryContributor, data)
}

func writeStoryFooter(buf *bytes.Buffer, ch git.Chronology, style styles.Style) error {
	if style.Templates.StoryFooter == "" {
		return nil
	}
	buf.WriteString("\n\n")
	return executeInline(buf, style.Templates.StoryFooter, buildStoryData(ch))
}
