package git

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

//go:embed .gitconfig
var configTemplate string

type gitClient struct {
	accessToken string
	owner       string
	repo        string
}

func NewGitClient(accessToken, owner string, repo string) (*gitClient, error) {
	f, err := os.Create(fmt.Sprintf("%s/.gitconfig", workDir))
	if err != nil {
		return nil, fmt.Errorf("failed to create .gitconfig file: %w", err)
	}

	// TODO(berkayoz): do not hardcode the name and user id
	config := map[string]string{
		"name":   "dev-ck8s-bot",
		"userid": "211096573",
	}

	tmpl, err := template.New("gitconfig").Parse(configTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(f, config); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close .gitconfig file: %w", err)
	}

	return &gitClient{
		accessToken: accessToken,
		owner:       owner,
		repo:        repo,
	}, nil
}

func (c *gitClient) GetRepositoryDir() string {
	return fmt.Sprintf("%s/%s", workDir, c.repo)
}

func (c *gitClient) GetRepositoryURL() string {
	return fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", c.accessToken, c.owner, c.repo)
}

func (c *gitClient) CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Dir = c.GetRepositoryDir()
	cmd.Env = append(cmd.Env, envOptions...)
	return cmd
}

func (c *gitClient) Clean(ctx context.Context) error {
	if err := os.RemoveAll(c.GetRepositoryDir()); err != nil {
		return fmt.Errorf("failed to remove repository directory: %w", err)
	}

	return nil
}
