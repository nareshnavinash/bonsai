package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/ui"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed models",
	RunE: func(cmd *cobra.Command, args []string) error {
		models, err := registry.ScanLocal()
		if err != nil {
			return fmt.Errorf("cannot scan models directory: %w", err)
		}

		if len(models) == 0 {
			fmt.Println("No models found. Pull a model with: bonsai pull <model>")
			return nil
		}

		headers := []string{"NAME", "SIZE", "MODIFIED"}
		rows := make([][]string, len(models))
		for i, m := range models {
			rows[i] = []string{
				m.Name,
				ui.FormatBytes(m.Size),
				ui.FormatRelativeTime(m.ModifiedAt),
			}
		}

		ui.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
