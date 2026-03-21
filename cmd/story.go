package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/narrative"
	"github.com/alciller88/commitlore/internal/renderer"
	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

const topPeaks = 3

func init() {
	rootCmd.AddCommand(newStoryCmd())
}

func newStoryCmd() *cobra.Command {
	var repoPath, from, style, format, output string

	cmd := &cobra.Command{
		Use:   "story",
		Short: "Generate a narrative of the repository history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStory(repoPath, from, style, format, output)
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Path to local git repository")
	cmd.Flags().StringVar(&from, "from", "", "Start from commit or tag (default: first commit)")
	cmd.Flags().StringVar(&style, "style", "formal", "Style: formal, patchnotes, ironic, epic")
	cmd.Flags().StringVar(&format, "format", "terminal", "Format: terminal, md, json")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (default: stdout)")

	return cmd
}

func runStory(repoPath, from, styleName, format, output string) error {
	repo, err := git.Open(repoPath)
	if err != nil {
		return err
	}

	commits, err := fetchStoryCommits(repo, from)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found in this repository.")
		return nil
	}

	ch, err := buildStoryChronology(repo, commits)
	if err != nil {
		return err
	}

	return renderStory(ch, styleName, format, output)
}

func fetchStoryCommits(repo *git.Repo, from string) ([]git.Commit, error) {
	opts := git.LogOptions{}
	if from != "" {
		t, err := time.Parse(dateLayout, from)
		if err != nil {
			return nil, fmt.Errorf("invalid --from value %q: use YYYY-MM-DD", from)
		}
		opts.Since = t
	}
	return repo.Log(opts)
}

func buildStoryChronology(repo *git.Repo, commits []git.Commit) (git.Chronology, error) {
	tags, err := repo.Tags()
	if err != nil {
		return git.Chronology{}, err
	}
	return git.BuildChronology(commits, tags, topPeaks), nil
}

func renderStory(ch git.Chronology, styleName, format, output string) error {
	style, err := styles.Load(styleName)
	if err != nil {
		return err
	}

	text, err := narrative.GenerateStory(ch, style)
	if err != nil {
		return err
	}

	rendered, err := renderer.RenderStory(text, ch, renderer.Format(format))
	if err != nil {
		return err
	}

	return writeOutput(rendered, output)
}
