package util

import (
	"context"
	"fmt"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v71/github"
)

func GetInstallationTokenFromClient(ctx context.Context, client *github.Client) (string, error) {
	itr := client.Client().Transport.(*ghinstallation.Transport)

	token, err := itr.Token(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get installation token: %w", err)
	}

	return token, nil
}
