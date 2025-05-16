package git

import (
	"context"
	"fmt"
)

func (c *gitClient) ForcePush(ctx context.Context) error {
	cmd := c.CommandContext(ctx, "git", "push", "--force")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	return nil
}
