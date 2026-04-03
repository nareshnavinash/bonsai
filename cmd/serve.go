package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
)

var serveCmd = &cobra.Command{
	Use:   "serve [model]",
	Short: "Start the llama-server",
	RunE: func(cmd *cobra.Command, args []string) error {
		model := defaultModel
		if len(args) > 0 {
			model = args[0]
		}

		modelPath, err := registry.ResolveModelPath(model)
		if err != nil {
			return err
		}

		mgr := getServerManager()
		mgr.Config.ModelPath = modelPath

		fmt.Printf("Starting llama-server with %s on port %d...\n", modelPath, mgr.Config.Port)
		fmt.Println("Press Ctrl+C to stop.")

		return mgr.StartForeground()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
