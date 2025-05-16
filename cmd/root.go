package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/berkayoz/maintainerd/internal"
	"github.com/spf13/cobra"
)

var githubAppID int64
var githubWebhookSecret string
var githubPrivateKeyFile string

var rootCmd = &cobra.Command{
	Use:   "maintainerd",
	Short: "A GitHub app that helps with repository and version maintenance",
	Long:  `Maintainerd is a GitHub app that helps with repository and version maintenance. It automates tasks such as updating dependencies, creating releases, and managing issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdCtx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill, syscall.SIGTERM)
		defer stop()

		mux := http.NewServeMux()

		eventMonitor, error := internal.NewGithubEventMonitor(githubAppID, githubWebhookSecret, githubPrivateKeyFile)
		if error != nil {
			return fmt.Errorf("failed to create GitHub event monitor: %w", error)
		}

		mux.Handle("/", eventMonitor)

		server := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}

		go func() {
			log.Print("Listening...")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP server ListenAndServe: %v", err)
			}
		}()

		<-cmdCtx.Done()

		log.Print("Shutting down server...")
		if err := server.Shutdown(cmdCtx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		return nil
	},
}

func Execute() {
	rootCmd.PersistentFlags().Int64Var(&githubAppID, "github-app-id", 0, "GitHub App ID")
	rootCmd.PersistentFlags().StringVar(&githubWebhookSecret, "github-webhook-secret", "", "GitHub Webhook Secret")
	rootCmd.PersistentFlags().StringVar(&githubPrivateKeyFile, "github-private-key-file", "", "GitHub Private Key File")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
