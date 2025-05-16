package git

import (
	"context"
	"fmt"
)

func (c *gitClient) Rebase(ctx context.Context, branch string) error {
	cmd := c.CommandContext(ctx, "git", "rebase", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebase branch: %w", err)
	}

	return nil
}
