package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	ghpkg "github.com/alciller88/commitlore/internal/github"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// GitApp exposes git history and contributor data to the frontend.
type GitApp struct{}

// NewGitApp creates a new GitApp instance.
func NewGitApp() *GitApp {
	return &GitApp{}
}

// OpenFolderPicker opens a native folder selection dialog.
// Returns the selected path or empty string if cancelled.
func (g *GitApp) OpenFolderPicker() string {
	app := application.Get()
	if app == nil {
		return ""
	}

	path, err := app.Dialog.OpenFile().
		CanChooseFiles(false).
		CanChooseDirectories(true).
		SetTitle("Select Repository Folder").
		PromptForSingleSelection()
	if err != nil {
		return ""
	}

	return path
}

// SetGitHubToken sets the GITHUB_TOKEN env var for the current session.
// The token is not persisted to disk. Pass empty string to clear.
func (g *GitApp) SetGitHubToken(token string) {
	os.Setenv("GITHUB_TOKEN", token)
}

// History returns commits as JSON, supporting local and GitHub repos.
func (g *GitApp) History(repo, author, since, until string, limit int) (string, error) {
	opts, err := buildOpts(author, since, until, limit)
	if err != nil {
		return "", fmt.Errorf("invalid filters: %w", err)
	}

	commits, err := fetchCommits(repo, opts)
	if err != nil {
		return "", cleanError(err)
	}

	return toJSON(commits)
}

// Contributors returns aggregated contributor data as JSON.
func (g *GitApp) Contributors(repo, since string, top int) (string, error) {
	opts := git.LogOptions{}
	if since != "" {
		t, err := time.Parse("2006-01-02", since)
		if err != nil {
			return "", fmt.Errorf("invalid since date: use YYYY-MM-DD")
		}
		opts.Since = t
	}

	commits, err := fetchCommits(repo, opts)
	if err != nil {
		return "", cleanError(err)
	}

	contribs := aggregateContribs(commits, top)
	return toJSON(contribs)
}

func buildOpts(author, since, until string, limit int) (git.LogOptions, error) {
	opts := git.LogOptions{Author: author, Limit: limit}

	if since != "" {
		t, err := time.Parse("2006-01-02", since)
		if err != nil {
			return opts, fmt.Errorf("invalid since date: use YYYY-MM-DD")
		}
		opts.Since = t
	}

	if until != "" {
		t, err := time.Parse("2006-01-02", until)
		if err != nil {
			return opts, fmt.Errorf("invalid until date: use YYYY-MM-DD")
		}
		opts.Until = t
	}

	return opts, nil
}

func fetchCommits(repoRef string, opts git.LogOptions) ([]git.Commit, error) {
	if ghpkg.IsRemoteRepo(repoRef) {
		owner, repo, err := ghpkg.ParseRepoRef(repoRef)
		if err != nil {
			return nil, err
		}
		client := ghpkg.NewClient(owner, repo)
		return client.ListCommits(context.Background(), opts)
	}

	r, err := git.Open(repoRef)
	if err != nil {
		return nil, err
	}
	return r.Log(opts)
}

type contribEntry struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Commits int    `json:"commits"`
}

func aggregateContribs(commits []git.Commit, top int) []contribEntry {
	counts := make(map[string]*contribEntry)
	for _, c := range commits {
		e, ok := counts[c.Email]
		if !ok {
			e = &contribEntry{Name: c.Author, Email: c.Email}
			counts[c.Email] = e
		}
		e.Commits++
	}

	sorted := sortContribs(counts, top)
	return sorted
}

func sortContribs(m map[string]*contribEntry, top int) []contribEntry {
	result := make([]contribEntry, 0, len(m))
	for _, e := range m {
		result = append(result, *e)
	}

	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Commits > result[i].Commits {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if top > 0 && len(result) > top {
		result = result[:top]
	}
	return result
}

func toJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("serialization error")
	}
	return string(data), nil
}

func cleanError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s", err.Error())
}
