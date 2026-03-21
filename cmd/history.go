package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/spf13/cobra"
)

const (
	defaultLimit    = 50
	hashDisplayLen  = 7
	dateLayout      = "2006-01-02"
	maxMessageWidth = 50
)

func init() {
	rootCmd.AddCommand(newHistoryCmd())
}

func newHistoryCmd() *cobra.Command {
	var repoPath, author, since, until string
	var limit int

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Explore commits filtered by author, date or range",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := buildLogOptions(author, since, until, limit)
			if err != nil {
				return err
			}
			return runHistory(repoPath, opts)
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Local path or GitHub repo (owner/repo)")
	cmd.Flags().StringVar(&author, "author", "", "Filter by author name or email")
	cmd.Flags().StringVar(&since, "since", "", "Show commits since date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "Show commits until date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&limit, "limit", defaultLimit, "Maximum number of commits")

	return cmd
}

func buildLogOptions(author, since, until string, limit int) (git.LogOptions, error) {
	opts := git.LogOptions{Author: author, Limit: limit}

	if since != "" {
		t, err := time.Parse(dateLayout, since)
		if err != nil {
			return opts, fmt.Errorf("invalid --since date %q: use YYYY-MM-DD", since)
		}
		opts.Since = t
	}

	if until != "" {
		t, err := time.Parse(dateLayout, until)
		if err != nil {
			return opts, fmt.Errorf("invalid --until date %q: use YYYY-MM-DD", until)
		}
		opts.Until = t
	}

	return opts, nil
}

func runHistory(repoPath string, opts git.LogOptions) error {
	commits, err := fetchCommitsFromSource(repoPath, opts)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found matching the given filters.")
		return nil
	}

	printCommitTable(commits)
	return nil
}

func printCommitTable(commits []git.Commit) {
	fmt.Printf("\n\033[1m%-9s %-12s %-20s %s\033[0m\n",
		"HASH", "DATE", "AUTHOR", "MESSAGE")
	fmt.Println("\033[90m" + repeatDash(80) + "\033[0m")

	for _, c := range commits {
		printCommitRow(c)
	}

	fmt.Printf("\n\033[90m%d commit(s)\033[0m\n", len(commits))
}

func printCommitRow(c git.Commit) {
	msg := c.Message
	if len(msg) > maxMessageWidth {
		msg = msg[:maxMessageWidth-3] + "..."
	}
	fmt.Printf("\033[33m%-9s\033[0m %-12s \033[36m%-20s\033[0m %s\n",
		c.Hash[:hashDisplayLen],
		c.Date.Format(dateLayout),
		truncate(c.Author, 20),
		msg)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func repeatDash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '-'
	}
	return string(b)
}
