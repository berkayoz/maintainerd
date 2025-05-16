package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/go-github/v71/github"
)

type githubEventHandler struct {
	client *github.Client
}

func New(client *github.Client) *githubEventHandler {
	return &githubEventHandler{
		client: client,
	}
}

func (h *githubEventHandler) Handle(ctx context.Context, event any) error {
	switch event := event.(type) {
	case *github.IssueCommentEvent:
		return h.processIssueCommentEvent(ctx, event)
	default:
		slog.Debug(fmt.Sprintf("unsupported event type: %T", event))
		return nil
	}
}
