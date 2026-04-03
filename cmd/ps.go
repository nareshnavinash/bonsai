package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/ui"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show running server status",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := getServerManager()
		info, err := mgr.ProcessInfo()
		if err != nil {
			fmt.Println("No server running.")
			return nil
		}

		modelName := registry.ModelFileName(info.ModelPath)
		// Try to map filename back to friendly name
		for _, m := range registry.Models {
			if m.GGUFFile == modelName {
				modelName = m.Name
				break
			}
		}

		headers := []string{"PID", "MODEL", "PORT", "UPTIME"}
		rows := [][]string{{
			fmt.Sprintf("%d", info.PID),
			modelName,
			fmt.Sprintf("%d", info.Port),
			ui.FormatRelativeTime(info.StartedAt),
		}}
		ui.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
