package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	ghpkg "github.com/alciller88/commitlore/internal/github"
	"github.com/spf13/cobra"
)

const defaultTop = 10

// Contributor holds aggregated stats for a single author.
type Contributor struct {
	Name       string
	Email      string
	Commits    int
	FirstDate  time.Time
	LastDate   time.Time
	FileCounts map[string]int
}

func init() {
	rootCmd.AddCommand(newContributorsCmd())
}

func newContributorsCmd() *cobra.Command {
	var repoPath, since string
	var top int

	cmd := &cobra.Command{
		Use:   "contributors",
		Short: "Show who has contributed and what they touched",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := buildContribOptions(since)
			if err != nil {
				return err
			}
			return runContributors(repoPath, opts, top)
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Local path or GitHub repo (owner/repo)")
	cmd.Flags().StringVar(&since, "since", "", "Analyze commits since date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&top, "top", defaultTop, "Number of contributors to show")

	return cmd
}

func buildContribOptions(since string) (git.LogOptions, error) {
	opts := git.LogOptions{}
	if since != "" {
		t, err := time.Parse(dateLayout, since)
		if err != nil {
			return opts, fmt.Errorf("invalid --since date %q: use YYYY-MM-DD", since)
		}
		opts.Since = t
	}
	return opts, nil
}

func runContributors(repoPath string, opts git.LogOptions, top int) error {
	commits, err := fetchCommitsFromSource(repoPath, opts)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found matching the given filters.")
		return nil
	}

	var contribs map[string]*Contributor
	if ghpkg.IsRemoteRepo(repoPath) {
		contribs = aggregateContributorsNoFiles(commits)
	} else {
		contribs = aggregateContributorsLocal(repoPath, commits)
	}

	ranked := rankByCommits(contribs, top)
	printContributorsTable(ranked)
	return nil
}

func aggregateContributorsLocal(repoPath string, commits []git.Commit) map[string]*Contributor {
	repo, err := git.Open(repoPath)
	if err != nil {
		return aggregateContributorsNoFiles(commits)
	}
	return aggregateContributors(repo, commits)
}

func aggregateContributors(repo *git.Repo, commits []git.Commit) map[string]*Contributor {
	contribs := make(map[string]*Contributor)
	for _, c := range commits {
		contrib := getOrCreateContributor(contribs, c)
		contrib.Commits++
		updateDateRange(contrib, c.Date)
		collectFiles(repo, contrib, c.Hash)
	}
	return contribs
}

func aggregateContributorsNoFiles(commits []git.Commit) map[string]*Contributor {
	contribs := make(map[string]*Contributor)
	for _, c := range commits {
		contrib := getOrCreateContributor(contribs, c)
		contrib.Commits++
		updateDateRange(contrib, c.Date)
	}
	return contribs
}

func getOrCreateContributor(m map[string]*Contributor, c git.Commit) *Contributor {
	if existing, ok := m[c.Email]; ok {
		return existing
	}
	contrib := &Contributor{
		Name:       c.Author,
		Email:      c.Email,
		FirstDate:  c.Date,
		LastDate:   c.Date,
		FileCounts: make(map[string]int),
	}
	m[c.Email] = contrib
	return contrib
}

func updateDateRange(c *Contributor, date time.Time) {
	if date.Before(c.FirstDate) {
		c.FirstDate = date
	}
	if date.After(c.LastDate) {
		c.LastDate = date
	}
}

func collectFiles(repo *git.Repo, contrib *Contributor, hash string) {
	files, err := repo.CommitFiles(hash)
	if err != nil {
		return
	}
	for _, f := range files {
		contrib.FileCounts[f]++
	}
}

func rankByCommits(contribs map[string]*Contributor, top int) []Contributor {
	sorted := make([]Contributor, 0, len(contribs))
	for _, c := range contribs {
		sorted = append(sorted, *c)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Commits > sorted[j].Commits
	})
	if top > 0 && len(sorted) > top {
		sorted = sorted[:top]
	}
	return sorted
}

func printContributorsTable(contribs []Contributor) {
	fmt.Printf("\n\033[1m%-25s %-8s %-25s %s\033[0m\n",
		"CONTRIBUTOR", "COMMITS", "ACTIVE PERIOD", "TOP FILES")
	fmt.Println("\033[90m" + repeatDash(90) + "\033[0m")

	for _, c := range contribs {
		printContributorRow(c)
	}
	fmt.Printf("\n\033[90m%d contributor(s)\033[0m\n", len(contribs))
}

func printContributorRow(c Contributor) {
	period := fmt.Sprintf("%s — %s",
		c.FirstDate.Format(dateLayout),
		c.LastDate.Format(dateLayout))

	topFiles := topNFiles(c.FileCounts, 3)

	fmt.Printf("\033[36m%-25s\033[0m %-8d %-25s %s\n",
		truncate(c.Name, 25),
		c.Commits,
		period,
		topFiles)
}

func topNFiles(counts map[string]int, n int) string {
	type fileCount struct {
		name  string
		count int
	}
	sorted := make([]fileCount, 0, len(counts))
	for name, count := range counts {
		sorted = append(sorted, fileCount{name, count})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})
	if len(sorted) > n {
		sorted = sorted[:n]
	}
	result := ""
	for i, fc := range sorted {
		if i > 0 {
			result += ", "
		}
		result += fc.name
	}
	return result
}
