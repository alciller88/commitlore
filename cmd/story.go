package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	ghpkg "github.com/alciller88/commitlore/internal/github"
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
	var llmFlag, llmBaseURL string

	cmd := &cobra.Command{
		Use:   "story",
		Short: "Generate a narrative of the repository history",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := resolveLLMProvider(llmFlag)
			baseURL := resolveLLMBaseURL(llmBaseURL)
			return runStory(repoPath, from, style, format, output, provider, baseURL)
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Local path or GitHub repo (owner/repo)")
	cmd.Flags().StringVar(&from, "from", "", "Start from commit or tag (default: first commit)")
	cmd.Flags().StringVar(&style, "style", "formal", "Style: formal, patchnotes, ironic, epic")
	cmd.Flags().StringVar(&format, "format", "terminal", "Format: terminal, md, json, html")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (default: stdout)")
	cmd.Flags().StringVar(&llmFlag, "llm", "none", "LLM provider: anthropic, openai, ollama, groq, none")
	cmd.Flags().StringVar(&llmBaseURL, "llm-base-url", "", "Override API base URL (OpenAI-compatible)")

	return cmd
}

func runStory(repoPath, from, styleName, format, output, llmProvider, llmBaseURL string) error {
	if ghpkg.IsRemoteRepo(repoPath) {
		return runRemoteStory(repoPath, from, styleName, format, output, llmProvider, llmBaseURL)
	}
	return runLocalStory(repoPath, from, styleName, format, output, llmProvider, llmBaseURL)
}

func runLocalStory(repoPath, from, styleName, format, output, llmProvider, llmBaseURL string) error {
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

	return renderStory(ch, styleName, format, output, llmProvider, llmBaseURL)
}

func runRemoteStory(repoRef, from, styleName, format, output, llmProvider, llmBaseURL string) error {
	opts := git.LogOptions{}
	if from != "" {
		t, err := time.Parse(dateLayout, from)
		if err != nil {
			return fmt.Errorf("invalid --from value %q: use YYYY-MM-DD", from)
		}
		opts.Since = t
	}

	commits, err := fetchRemoteCommits(repoRef, opts)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found in this repository.")
		return nil
	}

	ch := git.BuildChronology(commits, nil, topPeaks)
	return renderStory(ch, styleName, format, output, llmProvider, llmBaseURL)
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

func renderStory(ch git.Chronology, styleName, format, output, llmProvider, llmBaseURL string) error {
	style, err := styles.Load(styleName)
	if err != nil {
		return err
	}

	text, err := narrative.GenerateStory(ch, style)
	if err != nil {
		return err
	}

	text = enrichWithLLM(llmProvider, llmBaseURL, style.LLMPrompt, text)

	rendered, err := renderer.RenderStory(text, ch, style, renderer.Format(format))
	if err != nil {
		return err
	}

	return writeOutput(rendered, output)
}
