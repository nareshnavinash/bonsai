package cmd

import (
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/ui"
)

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available Bonsai models from HuggingFace",
	RunE: func(cmd *cobra.Command, args []string) error {
		headers := []string{"NAME", "PARAMS", "SIZE", "PULL COMMAND"}
		rows := make([][]string, len(registry.Models))
		for i, m := range registry.Models {
			rows[i] = []string{
				m.Name,
				m.Params,
				m.Size,
				"bonsai pull " + m.Name,
			}
		}

		ui.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}
