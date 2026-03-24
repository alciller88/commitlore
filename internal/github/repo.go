package github

import (
	"context"
	"strings"

	"github.com/alciller88/commitlore/internal/git"
	gh "github.com/google/go-github/v62/github"
)

// ListCommits fetches commits from a GitHub repo applying the same filters
// as internal/git: since, until, author, limit.
func (c *Client) ListCommits(ctx context.Context, opts git.LogOptions) ([]git.Commit, error) {
	ghOpts := c.buildCommitListOptions(opts)
	return c.fetchAllCommits(ctx, ghOpts, opts)
}

func (c *Client) buildCommitListOptions(opts git.LogOptions) *gh.CommitsListOptions {
	ghOpts := &gh.CommitsListOptions{
		ListOptions: gh.ListOptions{PerPage: perPageDefault},
	}

	if opts.Author != "" {
		ghOpts.Author = opts.Author
	}
	if !opts.Since.IsZero() {
		ghOpts.Since = opts.Since
	}
	if !opts.Until.IsZero() {
		ghOpts.Until = opts.Until
	}

	return ghOpts
}

const perPageDefault = 100

func (c *Client) fetchAllCommits(
	ctx context.Context,
	ghOpts *gh.CommitsListOptions,
	opts git.LogOptions,
) ([]git.Commit, error) {
	var commits []git.Commit

	for {
		batch, resp, err := c.gh.Repositories.ListCommits(ctx, c.owner, c.repo, ghOpts)
		if err != nil {
			return nil, wrapAPIError(err, c.owner, c.repo)
		}

		for _, rc := range batch {
			commit := toGitCommit(rc)
			if !matchesAuthorFilter(commit, opts.Author) {
				continue
			}
			commits = append(commits, commit)
			if opts.Limit > 0 && len(commits) >= opts.Limit {
				return commits, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}
		ghOpts.Page = resp.NextPage
	}

	return commits, nil
}

func toGitCommit(rc *gh.RepositoryCommit) git.Commit {
	c := git.Commit{
		Hash: rc.GetSHA(),
	}

	if gc := rc.GetCommit(); gc != nil {
		c.Message = subjectLine(gc.GetMessage())

		if author := gc.GetAuthor(); author != nil {
			c.Author = author.GetName()
			c.Email = author.GetEmail()
			c.Date = author.GetDate().Time
		}
	}

	return c
}

// matchesAuthorFilter applies substring matching on author name/email,
// complementing the GitHub API's author filter (which only matches login).
func matchesAuthorFilter(c git.Commit, author string) bool {
	if author == "" {
		return true
	}
	lower := strings.ToLower(author)
	return strings.Contains(strings.ToLower(c.Author), lower) ||
		strings.Contains(strings.ToLower(c.Email), lower)
}

func subjectLine(msg string) string {
	msg = strings.TrimSpace(msg)
	if idx := strings.IndexByte(msg, '\n'); idx >= 0 {
		msg = msg[:idx]
	}
	return strings.TrimSpace(msg)
}
