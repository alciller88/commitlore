package app

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/git"
	ghpkg "github.com/alciller88/commitlore/internal/github"
	"github.com/alciller88/commitlore/internal/narrative"
	"github.com/alciller88/commitlore/internal/renderer"
	"github.com/alciller88/commitlore/internal/styles"
)

const storyTopPeaks = 3

// StoryApp exposes story generation to the frontend.
type StoryApp struct{}

// NewStoryApp creates a new StoryApp instance.
func NewStoryApp() *StoryApp {
	return &StoryApp{}
}

// GenerateStory produces a narrative as HTML from the given repo.
func (s *StoryApp) GenerateStory(repo, from, styleName, llmProvider, llmModel string) (string, error) {
	opts, err := buildOpts("", from, "", 0)
	if err != nil {
		return "", err
	}

	commits, err := fetchCommits(repo, opts)
	if err != nil {
		return "", cleanError(err)
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found in this repository")
	}

	ch, err := buildChronology(repo, commits)
	if err != nil {
		return "", cleanError(err)
	}

	return renderStoryHTML(ch, styleName, llmProvider, llmModel)
}

func buildChronology(repoRef string, commits []git.Commit) (git.Chronology, error) {
	if ghpkg.IsRemoteRepo(repoRef) {
		return git.BuildChronology(commits, nil, storyTopPeaks), nil
	}

	r, err := git.Open(repoRef)
	if err != nil {
		return git.Chronology{}, cleanError(err)
	}

	tags, err := r.Tags()
	if err != nil {
		return git.Chronology{}, cleanError(err)
	}

	return git.BuildChronology(commits, tags, storyTopPeaks), nil
}

func renderStoryHTML(ch git.Chronology, styleName, llmProvider, llmModel string) (string, error) {
	style, err := styles.Load(styleName)
	if err != nil {
		return "", cleanError(err)
	}

	text, err := narrative.GenerateStory(ch, style)
	if err != nil {
		return "", cleanError(err)
	}

	text = tryEnrich(llmProvider, llmModel, style.LLMPrompt, text)

	rendered, err := renderer.RenderStory(text, ch, style, renderer.Format("html"))
	if err != nil {
		return "", cleanError(err)
	}

	return rendered, nil
}
