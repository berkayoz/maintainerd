package git

import (
	"context"
	"fmt"
	"log/slog"
)

func (c *gitClient) Clone(ctx context.Context) error {
	repoDir := c.GetRepositoryDir()
	url := c.GetRepositoryURL()

	cmd := c.CommandContext(ctx, "git", "clone", url, repoDir)
	cmd.Dir = workDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to clone repository", "error", err, "output", string(out))
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}
