package cmd

import (
	"fmt"
	"os"

	"github.com/alciller88/commitlore/internal/changelog"
	"github.com/alciller88/commitlore/internal/git"
	"github.com/alciller88/commitlore/internal/narrative"
	"github.com/alciller88/commitlore/internal/renderer"
	"github.com/alciller88/commitlore/internal/styles"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newGenerateCmd())
}

func newGenerateCmd() *cobra.Command {
	var repoPath, since, until, style, format, output string
	var includePRs bool

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a changelog from commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			if includePRs {
				// TODO(fase8): implement PR inclusion via GitHub API
				fmt.Fprintln(os.Stderr, "Warning: --include-prs is not yet implemented, ignoring flag.")
			}
			return runGenerate(repoPath, since, until, style, format, output)
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Path to local git repository")
	cmd.Flags().StringVar(&since, "since", "", "Show commits since date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "Show commits until date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&includePRs, "include-prs", false, "Include PR info (requires GitHub token)")
	cmd.Flags().StringVar(&style, "style", "formal", "Style: formal, patchnotes, ironic, epic")
	cmd.Flags().StringVar(&format, "format", "terminal", "Format: terminal, md, json, html, pdf")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (default: stdout)")

	return cmd
}

func runGenerate(repoPath, since, until, styleName, format, output string) error {
	opts, err := buildLogOptions("", since, until, 0)
	if err != nil {
		return err
	}

	commits, err := fetchCommits(repoPath, opts)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found matching the given filters.")
		return nil
	}

	cl := changelog.GroupCommits(commits)
	return generateOutput(cl, styleName, format, output)
}

func fetchCommits(repoPath string, opts git.LogOptions) ([]git.Commit, error) {
	repo, err := git.Open(repoPath)
	if err != nil {
		return nil, err
	}
	return repo.Log(opts)
}

func generateOutput(cl changelog.Changelog, styleName, format, output string) error {
	style, err := styles.Load(styleName)
	if err != nil {
		return err
	}

	text, err := narrative.Generate(cl, style)
	if err != nil {
		return err
	}

	rendered, err := renderer.Render(text, cl, renderer.Format(format))
	if err != nil {
		return err
	}

	return writeOutput(rendered, output)
}

func writeOutput(content, output string) error {
	if output == "" {
		fmt.Print(content)
		return nil
	}
	err := os.WriteFile(output, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("writing to %s: %w", output, err)
	}
	fmt.Fprintf(os.Stderr, "Output written to %s\n", output)
	return nil
}
