package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alciller88/commitlore/internal/git"
	ghpkg "github.com/alciller88/commitlore/internal/github"
)

// fetchCommitsFromSource fetches commits from a local repo or GitHub,
// depending on the format of repoRef.
func fetchCommitsFromSource(repoRef string, opts git.LogOptions) ([]git.Commit, error) {
	if ghpkg.IsRemoteRepo(repoRef) {
		return fetchRemoteCommits(repoRef, opts)
	}
	return fetchLocalCommits(repoRef, opts)
}

func fetchLocalCommits(repoPath string, opts git.LogOptions) ([]git.Commit, error) {
	repo, err := git.Open(repoPath)
	if err != nil {
		return nil, err
	}
	return repo.Log(opts)
}

func fetchRemoteCommits(repoRef string, opts git.LogOptions) ([]git.Commit, error) {
	owner, repo, err := ghpkg.ParseRepoRef(repoRef)
	if err != nil {
		return nil, err
	}
	client := ghpkg.NewClient(owner, repo)
	return client.ListCommits(context.Background(), opts)
}

// fetchRemotePRs fetches merged PRs from a GitHub repo.
// Returns an error if no GITHUB_TOKEN is set.
func fetchRemotePRs(repoRef string, since, until time.Time) ([]ghpkg.PullRequest, error) {
	if !ghpkg.HasToken() {
		fmt.Fprintln(os.Stderr, "Warning: --include-prs requires GITHUB_TOKEN to be set. Skipping PR data.")
		return nil, nil
	}

	owner, repo, err := ghpkg.ParseRepoRef(repoRef)
	if err != nil {
		return nil, err
	}

	client := ghpkg.NewClient(owner, repo)
	return client.ListPullRequests(context.Background(), since, until)
}
