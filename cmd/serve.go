package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Ollama server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting Ollama server...")

		c := exec.Command("ollama", "serve")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to start Ollama: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
