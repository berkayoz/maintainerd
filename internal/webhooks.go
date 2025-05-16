package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/berkayoz/maintainerd/internal/handler"
	"github.com/berkayoz/maintainerd/internal/util"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v71/github"
)

type githubEventMonitor struct {
	webhookSecretKey []byte
	appID            int64
	privateKey       []byte
}

func NewGithubEventMonitor(appID int64, webhookSecretKey string, privateKeyFile string) (*githubEventMonitor, error) {
	privateKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %s", err)
	}

	return &githubEventMonitor{
		appID:            appID,
		webhookSecretKey: []byte(webhookSecretKey),
		privateKey:       privateKey,
	}, nil
}

// TODO(berkayoz): Maybe ignore errors and respon d with 200 OK for all webhook events
func (s *githubEventMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, s.webhookSecretKey)
	if err != nil {
		slog.Error("failed to validate payload", "error", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		slog.Error("failed to parse webhook", "error", err)
		http.Error(w, "invalid webhook", http.StatusBadRequest)
		return
	}

	eventMeta, err := util.UnmarshalEventMeta(payload)
	if err != nil {
		slog.Error("failed to unmarshal event meta", "error", err)
		http.Error(w, "invalid event meta", http.StatusBadRequest)
		return
	}

	installationID := eventMeta.Installation.GetID()

	client, err := s.newGithubClient(installationID)
	if err != nil {
		slog.Error("failed to create GitHub client", "error", err)
		http.Error(w, "failed to create GitHub client", http.StatusInternalServerError)
		return
	}

	eventHandler := handler.New(client)

	if err := eventHandler.Handle(r.Context(), event); err != nil {
		slog.Error("failed to handle event", "error", err)
		http.Error(w, "failed to handle event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func (s *githubEventMonitor) newGithubClient(installationID int64) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, s.appID, installationID, s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	// Use installation transport with github.com/google/go-github
	return github.NewClient(&http.Client{Transport: itr}), nil
}
