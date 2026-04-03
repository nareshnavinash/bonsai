package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the llama-server",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := getServerManager()
		if !mgr.IsRunning() {
			fmt.Println("Server is not running.")
			return nil
		}
		if err := mgr.Stop(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
		fmt.Println("Server stopped.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
