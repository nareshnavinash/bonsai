package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show server status",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := getServerManager()
		client := getLLMClient()

		host := client.BaseURL

		processRunning := mgr.IsRunning()
		serverHealthy := client.Health(context.Background()) == nil

		status := "stopped"
		if processRunning && serverHealthy {
			status = "running"
		} else if processRunning {
			status = "starting"
		}

		modelInfo := defaultModel + " (configured)"
		if processRunning {
			modelPath := mgr.LoadedModel()
			modelName := registry.ModelFileName(modelPath)
			for _, m := range registry.Models {
				if m.GGUFFile == modelName {
					modelName = m.Name
					break
				}
			}
			modelInfo = modelName + " (loaded)"
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
