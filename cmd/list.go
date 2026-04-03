package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshsekar/bonsai/internal/ui"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available models",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		resp, err := client.List(context.Background())
		if err != nil {
			return fmt.Errorf("cannot connect to Ollama: %w", err)
		}

		if len(resp.Models) == 0 {
			fmt.Println("No models found. Pull a model with: bonsai pull <model>")
			return nil
		}

		headers := []string{"NAME", "ID", "SIZE", "MODIFIED"}
		rows := make([][]string, len(resp.Models))
		for i, m := range resp.Models {
			rows[i] = []string{
				m.Name,
				ui.TruncateID(m.Digest),
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
