package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/berkayoz/maintainerd/internal/git"
	"github.com/berkayoz/maintainerd/internal/util"
	"github.com/google/go-github/v71/github"
)

func (h *githubEventHandler) processIssueCommentEvent(ctx context.Context, event *github.IssueCommentEvent) error {
	if event.GetAction() == "created" {
		comment := event.GetComment()
		slog.Info("New comment on issue", "issue_number", event.GetIssue().GetNumber(), "comment_body", comment.GetBody())

		switch strings.TrimSpace(comment.GetBody()) {
		case "/test":
			owner := event.GetRepo().GetOwner().GetLogin()
			repo := event.GetRepo().GetName()
			issueNumber := event.GetIssue().GetNumber()

			_, _, err := h.client.Issues.CreateComment(ctx, owner, repo, issueNumber, &github.IssueComment{
				Body: github.Ptr("Thank you for your comment!"),
			})

			if err != nil {
				return fmt.Errorf("failed to create comment: %w", err)
			}
		case "/rebase":
			if event.GetIssue().IsPullRequest() {
				slog.Info("The issue is a pull request", "issue_number", event.GetIssue().GetNumber())
			} else {
				slog.Info("The issue is not a pull request", "issue_number", event.GetIssue().GetNumber())
				return nil
			}

			owner := event.GetRepo().GetOwner().GetLogin()
			repo := event.GetRepo().GetName()
			issueNumber := event.GetIssue().GetNumber()

			pr, _, err := h.client.PullRequests.Get(ctx, owner, repo, issueNumber)
			if err != nil {
				return fmt.Errorf("failed to get pull request details: %w", err)
			}

			slog.Info("Pull request details retrieved", "title", pr.GetTitle(), "head", pr.GetHead().GetRef(), "base", pr.GetBase().GetRef())

			token, err := util.GetInstallationTokenFromClient(ctx, h.client)
			if err != nil {
				return fmt.Errorf("failed to get installation token: %w", err)
			}

			gitClient, err := git.NewGitClient(token, owner, repo)
			if err != nil {
				return fmt.Errorf("failed to create git client: %w", err)
			}
			defer func() {
				if err := gitClient.Clean(ctx); err != nil {
					slog.Error("Failed to clean up repository", "error", err)
				}
			}()

			if err := gitClient.Clone(ctx); err != nil {
				return fmt.Errorf("failed to clone repository: %w", err)
			}

			if err := gitClient.Checkout(ctx, pr.GetBase().GetRef()); err != nil {
				return fmt.Errorf("failed to checkout base branch: %w", err)
			}

			if err := gitClient.Checkout(ctx, pr.GetHead().GetRef()); err != nil {
				return fmt.Errorf("failed to checkout head branch: %w", err)
			}

			if err := gitClient.Rebase(ctx, pr.GetBase().GetRef()); err != nil {
				return fmt.Errorf("failed to rebase branch: %w", err)
			}

			if err := gitClient.ForcePush(ctx); err != nil {
				return fmt.Errorf("failed to push changes: %w", err)
			}

		}

	}

	return nil
}
