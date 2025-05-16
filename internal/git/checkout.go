package git

import (
	"context"
	"fmt"
)

func (c *gitClient) Checkout(ctx context.Context, branch string) error {
	cmd := c.CommandContext(ctx, "git", "checkout", branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	return nil
}
