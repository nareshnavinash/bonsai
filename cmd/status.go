package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show server status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.ClientFromEnvironment()
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

		fmt.Printf("Server:   %s\n", host)
		fmt.Printf("Model:    %s\n", defaultModel)
		fmt.Printf("Status:   %s\n", status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
