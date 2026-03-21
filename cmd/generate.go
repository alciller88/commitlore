package cmd

import (
	"fmt"
	"os"

	"github.com/alciller88/commitlore/internal/changelog"
	ghpkg "github.com/alciller88/commitlore/internal/github"
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
	var llmFlag, llmBaseURL string
	var includePRs bool

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a changelog from commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := resolveLLMProvider(llmFlag)
			baseURL := resolveLLMBaseURL(llmBaseURL)
			return runGenerate(repoPath, since, until, style, format, output, provider, baseURL, includePRs)
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Local path or GitHub repo (owner/repo)")
	cmd.Flags().StringVar(&since, "since", "", "Show commits since date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "Show commits until date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&includePRs, "include-prs", false, "Include PR info (requires GitHub token)")
	cmd.Flags().StringVar(&style, "style", "formal", "Style: formal, patchnotes, ironic, epic")
	cmd.Flags().StringVar(&format, "format", "terminal", "Format: terminal, md, json, html")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (default: stdout)")
	cmd.Flags().StringVar(&llmFlag, "llm", "none", "LLM provider: anthropic, openai, ollama, groq, none")
	cmd.Flags().StringVar(&llmBaseURL, "llm-base-url", "", "Override API base URL (OpenAI-compatible)")

	return cmd
}

func runGenerate(repoPath, since, until, styleName, format, output, llmProvider, llmBaseURL string, includePRs bool) error {
	opts, err := buildLogOptions("", since, until, 0)
	if err != nil {
		return err
	}

	commits, err := fetchCommitsFromSource(repoPath, opts)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Fprintln(os.Stderr, "No commits found matching the given filters.")
		return nil
	}

	cl := changelog.GroupCommits(commits)

	if includePRs && ghpkg.IsRemoteRepo(repoPath) {
		prs, prErr := fetchRemotePRs(repoPath, opts.Since, opts.Until)
		if prErr != nil {
			return prErr
		}
		if len(prs) > 0 {
			appendPRsToChangelog(&cl, prs)
		}
	} else if includePRs && !ghpkg.IsRemoteRepo(repoPath) {
		fmt.Fprintln(os.Stderr, "Warning: --include-prs requires --repo in owner/repo format. Ignoring flag.")
	}

	return generateOutput(cl, styleName, format, output, llmProvider, llmBaseURL)
}

func appendPRsToChangelog(cl *changelog.Changelog, prs []ghpkg.PullRequest) {
	for _, pr := range prs {
		msg := fmt.Sprintf("PR #%d: %s (by @%s)", pr.Number, pr.Title, pr.Author)
		cl.AppendCommit(changelog.ParsedCommit{
			Hash:    fmt.Sprintf("pr-%d", pr.Number),
			Message: msg,
			Type:    changelog.InferTypeFromMessage(msg),
		})
	}
}

func generateOutput(cl changelog.Changelog, styleName, format, output, llmProvider, llmBaseURL string) error {
	style, err := styles.Load(styleName)
	if err != nil {
		return err
	}

	text, err := narrative.Generate(cl, style)
	if err != nil {
		return err
	}

	text = enrichWithLLM(llmProvider, llmBaseURL, style.LLMPrompt, text)

	rendered, err := renderer.Render(text, cl, style, renderer.Format(format))
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
