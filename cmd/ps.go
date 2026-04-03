package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/ui"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running models",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		resp, err := client.ListRunning(context.Background())
		if err != nil {
			return fmt.Errorf("cannot connect to Ollama: %w", err)
		}

		if len(resp.Models) == 0 {
			fmt.Println("No models currently running.")
			return nil
		}

		headers := []string{"NAME", "ID", "SIZE", "PROCESSOR", "UNTIL"}
		rows := make([][]string, len(resp.Models))
		for i, m := range resp.Models {
			vramPct := 0
			if m.Size > 0 {
				vramPct = int(float64(m.SizeVRAM) / float64(m.Size) * 100)
			}
			processor := "100% CPU"
			if vramPct == 100 {
				processor = "100% GPU"
			} else if vramPct > 0 {
				processor = fmt.Sprintf("%d%% GPU / %d%% CPU", vramPct, 100-vramPct)
			}

			until := ui.FormatRelativeFuture(m.ExpiresAt)
			if m.ExpiresAt.Before(time.Now()) {
				until = "expired"
			}

			rows[i] = []string{
				m.Name,
				ui.TruncateID(m.Digest),
				ui.FormatBytes(int64(m.Size)),
				processor,
				until,
			}
		}

		ui.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
