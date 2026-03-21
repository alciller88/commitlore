package github

import (
	"context"
	"time"

	gh "github.com/google/go-github/v62/github"
)

// PullRequest represents a GitHub pull request (read-only data).
type PullRequest struct {
	Number    int
	Title     string
	Author    string
	State     string
	CreatedAt time.Time
	MergedAt  time.Time
	Labels    []string
}

// ListPullRequests fetches merged PRs from the repository.
// Only GET operations — never creates, updates, or deletes PRs.
func (c *Client) ListPullRequests(ctx context.Context, since, until time.Time) ([]PullRequest, error) {
	ghOpts := &gh.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: gh.ListOptions{
			PerPage: perPageDefault,
		},
	}

	return c.fetchAllPRs(ctx, ghOpts, since, until)
}

func (c *Client) fetchAllPRs(
	ctx context.Context,
	ghOpts *gh.PullRequestListOptions,
	since, until time.Time,
) ([]PullRequest, error) {
	var prs []PullRequest

	for {
		batch, resp, err := c.gh.PullRequests.List(ctx, c.owner, c.repo, ghOpts)
		if err != nil {
			return nil, wrapAPIError(err)
		}

		done := false
		for _, pr := range batch {
			if !isMerged(pr) {
				continue
			}
			merged := pr.GetMergedAt().Time
			if !since.IsZero() && merged.Before(since) {
				done = true
				break
			}
			if !until.IsZero() && merged.After(until) {
				continue
			}
			prs = append(prs, toPullRequest(pr))
		}

		if done || resp.NextPage == 0 {
			break
		}
		ghOpts.Page = resp.NextPage
	}

	return prs, nil
}

func isMerged(pr *gh.PullRequest) bool {
	return pr.MergedAt != nil && !pr.GetMergedAt().IsZero()
}

func toPullRequest(pr *gh.PullRequest) PullRequest {
	result := PullRequest{
		Number:    pr.GetNumber(),
		Title:     pr.GetTitle(),
		State:     pr.GetState(),
		CreatedAt: pr.GetCreatedAt().Time,
		MergedAt:  pr.GetMergedAt().Time,
	}

	if user := pr.GetUser(); user != nil {
		result.Author = user.GetLogin()
	}

	for _, label := range pr.Labels {
		result.Labels = append(result.Labels, label.GetName())
	}

	return result
}
