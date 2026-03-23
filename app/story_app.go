package main

import (
	"fmt"

	"github.com/alciller88/commitlore/internal/git"
	ghpkg "github.com/alciller88/commitlore/internal/github"
	"github.com/alciller88/commitlore/internal/narrative"
	"github.com/alciller88/commitlore/internal/renderer"
)

const storyTopPeaks = 12

// StoryApp exposes story generation to the frontend.
type StoryApp struct{}

// NewStoryApp creates a new StoryApp instance.
func NewStoryApp() *StoryApp {
	return &StoryApp{}
}

// GenerateStory produces a narrative as HTML from the given repo.
// LLM config is read automatically from Settings (config.yml + keychain).
func (s *StoryApp) GenerateStory(repo, from, styleName string) (string, error) {
	opts, err := buildOpts("", from, "", 0)
	if err != nil {
		return "", err
	}

	var commits []git.Commit
	if shouldUseCache(opts) {
		commits, _ = globalCommitCache.get(repo)
	}
	if commits == nil {
		commits, err = fetchCommits(repo, opts)
		if err != nil {
			return "", cleanError(err)
		}
		if shouldUseCache(opts) {
			globalCommitCache.set(repo, commits)
		}
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found in this repository")
	}

	ch, err := buildChronology(repo, commits)
	if err != nil {
		return "", cleanError(err)
	}

	provider, model := loadLLMSettings()
	repoName := renderer.RepoNameFromPath(repo)
	return renderStoryHTML(ch, styleName, provider, model, repoName)
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

func renderStoryHTML(ch git.Chronology, styleName, llmProvider, llmModel, repoName string) (string, error) {
	style, err := loadStyleWithLanguage(styleName)
	if err != nil {
		return "", cleanError(err)
	}

	text, err := narrative.GenerateStory(ch, style)
	if err != nil {
		return "", cleanError(err)
	}

	text = tryEnrich(llmProvider, llmModel, style.LLMPrompt, text)

	override := buildHTMLThemeOverride(styleName)
	rendered, err := renderer.RenderStoryWithTheme(text, ch, style, renderer.Format("html"), override, repoName)
	if err != nil {
		return "", cleanError(err)
	}

	return rendered, nil
}
