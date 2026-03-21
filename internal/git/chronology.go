package git

import (
	"sort"
	"time"
)

// Tag represents a git tag with its associated commit and date.
type Tag struct {
	Name string
	Hash string
	Date time.Time
}

// ActivityPeak represents a month with high commit activity.
type ActivityPeak struct {
	Month string // "2006-01"
	Count int
}

// ContributorEntry represents the first appearance of a contributor.
type ContributorEntry struct {
	Name  string
	Email string
	Date  time.Time
}

// Chronology holds the timeline data for a repository story.
type Chronology struct {
	FirstCommit  Commit
	Tags         []Tag
	Peaks        []ActivityPeak
	Contributors []ContributorEntry
	TotalCommits int
}

// BuildChronology constructs a Chronology from commits and tags.
func BuildChronology(commits []Commit, tags []Tag, topPeaks int) Chronology {
	if len(commits) == 0 {
		return Chronology{}
	}

	sorted := sortByDateAsc(commits)

	return Chronology{
		FirstCommit:  sorted[0],
		Tags:         tags,
		Peaks:        activityPeaks(sorted, topPeaks),
		Contributors: uniqueContributors(sorted),
		TotalCommits: len(sorted),
	}
}

func sortByDateAsc(commits []Commit) []Commit {
	result := make([]Commit, len(commits))
	copy(result, commits)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})
	return result
}

// activityPeaks returns the top N months by commit count, descending.
func activityPeaks(commits []Commit, top int) []ActivityPeak {
	counts := countByMonth(commits)
	peaks := mapToPeaks(counts)
	sort.Slice(peaks, func(i, j int) bool {
		return peaks[i].Count > peaks[j].Count
	})
	if top > 0 && len(peaks) > top {
		peaks = peaks[:top]
	}
	return peaks
}

func countByMonth(commits []Commit) map[string]int {
	counts := make(map[string]int)
	for _, c := range commits {
		key := c.Date.Format("2006-01")
		counts[key]++
	}
	return counts
}

func mapToPeaks(counts map[string]int) []ActivityPeak {
	peaks := make([]ActivityPeak, 0, len(counts))
	for month, count := range counts {
		peaks = append(peaks, ActivityPeak{Month: month, Count: count})
	}
	return peaks
}

// uniqueContributors returns contributors ordered by first appearance.
func uniqueContributors(commits []Commit) []ContributorEntry {
	seen := make(map[string]bool)
	var contributors []ContributorEntry

	for _, c := range commits {
		if seen[c.Email] {
			continue
		}
		seen[c.Email] = true
		contributors = append(contributors, ContributorEntry{
			Name:  c.Author,
			Email: c.Email,
			Date:  c.Date,
		})
	}
	return contributors
}
