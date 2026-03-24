package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildChronology_empty(t *testing.T) {
	ch := BuildChronology(nil, nil, 3)
	assert.Equal(t, 0, ch.TotalCommits)
	assert.Empty(t, ch.Peaks)
	assert.Empty(t, ch.Contributors)
}

func TestBuildChronology_firstCommit(t *testing.T) {
	commits := chronologyCommits()
	ch := BuildChronology(commits, nil, 3)
	assert.Equal(t, "first commit", ch.FirstCommit.Message)
	assert.Equal(t, 6, ch.TotalCommits)
}

func TestBuildChronology_preservesTags(t *testing.T) {
	tags := []Tag{
		{Name: "v1.0.0", Hash: "abc", Date: time.Now()},
	}
	ch := BuildChronology(chronologyCommits(), tags, 3)
	assert.Len(t, ch.Tags, 1)
	assert.Equal(t, "v1.0.0", ch.Tags[0].Name)
}

func TestActivityPeaks_topN(t *testing.T) {
	commits := chronologyCommits()
	peaks := activityPeaks(commits, 2)
	assert.Len(t, peaks, 2)
	assert.GreaterOrEqual(t, peaks[0].Count, peaks[1].Count)
}

func TestActivityPeaks_allMonths(t *testing.T) {
	commits := chronologyCommits()
	peaks := activityPeaks(commits, 0)
	assert.Len(t, peaks, 3)
}

func TestActivityPeaks_descendingOrder(t *testing.T) {
	commits := chronologyCommits()
	peaks := activityPeaks(commits, 3)
	for i := 1; i < len(peaks); i++ {
		assert.GreaterOrEqual(t, peaks[i-1].Count, peaks[i].Count)
	}
}

func TestUniqueContributors_order(t *testing.T) {
	commits := chronologyCommits()
	sorted := sortByDateAsc(commits)
	contributors := uniqueContributors(sorted)
	assert.Len(t, contributors, 3)
	assert.Equal(t, "Alice", contributors[0].Name)
	assert.Equal(t, "Bob", contributors[1].Name)
	assert.Equal(t, "Charlie", contributors[2].Name)
}

func TestUniqueContributors_noDuplicates(t *testing.T) {
	commits := chronologyCommits()
	sorted := sortByDateAsc(commits)
	contributors := uniqueContributors(sorted)
	emails := make(map[string]bool)
	for _, c := range contributors {
		assert.False(t, emails[c.Email], "duplicate: %s", c.Email)
		emails[c.Email] = true
	}
}

func TestSortByDateAsc(t *testing.T) {
	commits := chronologyCommits()
	sorted := sortByDateAsc(commits)
	for i := 1; i < len(sorted); i++ {
		assert.False(t, sorted[i].Date.Before(sorted[i-1].Date))
	}
}

func TestCountByMonth(t *testing.T) {
	commits := chronologyCommits()
	counts := countByMonth(commits)
	assert.Equal(t, 3, counts["2025-01"])
	assert.Equal(t, 2, counts["2025-02"])
	assert.Equal(t, 1, counts["2025-03"])
}

// chronologyCommits returns commits spread across 3 months with 3 authors.
func chronologyCommits() []Commit {
	return []Commit{
		{Hash: "aaa", Author: "Alice", Email: "alice@x.com", Date: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC), Message: "first commit"},
		{Hash: "bbb", Author: "Alice", Email: "alice@x.com", Date: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC), Message: "second commit"},
		{Hash: "ccc", Author: "Bob", Email: "bob@x.com", Date: time.Date(2025, 1, 20, 10, 0, 0, 0, time.UTC), Message: "third commit"},
		{Hash: "ddd", Author: "Bob", Email: "bob@x.com", Date: time.Date(2025, 2, 5, 10, 0, 0, 0, time.UTC), Message: "fourth commit"},
		{Hash: "eee", Author: "Charlie", Email: "charlie@x.com", Date: time.Date(2025, 2, 10, 10, 0, 0, 0, time.UTC), Message: "fifth commit"},
		{Hash: "fff", Author: "Alice", Email: "alice@x.com", Date: time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC), Message: "sixth commit"},
	}
}
