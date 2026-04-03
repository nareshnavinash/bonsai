package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show server status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		host := os.Getenv("OLLAMA_HOST")
		if host == "" {
			host = "http://localhost:11434"
		}

		status := "stopped"
		if err := client.Heartbeat(context.Background()); err == nil {
			status = "running"
		}

		// Show actual running model, not just configured default
		modelInfo := defaultModel + " (configured)"
		if resp, err := client.ListRunning(context.Background()); err == nil && len(resp.Models) > 0 {
			modelInfo = resp.Models[0].Name + " (loaded)"
		}

		fmt.Printf("Server:   %s\n", host)
		fmt.Printf("Model:    %s\n", modelInfo)
		fmt.Printf("Status:   %s\n", status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
