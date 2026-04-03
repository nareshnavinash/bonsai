package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
)

var rmCmd = &cobra.Command{
	Use:   "rm <model>",
	Short: "Remove a model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path, err := registry.ResolveModelPath(name)
		if err != nil {
			return fmt.Errorf("model %q not found locally: %w", name, err)
		}

		if err := os.Remove(path); err != nil {
			return fmt.Errorf("failed to delete %q: %w", path, err)
		}

		fmt.Printf("Deleted %s (%s)\n", name, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
