package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp <source> <destination>",
	Short: "Copy a model",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		source, dest := args[0], args[1]

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		err = client.Copy(context.Background(), &api.CopyRequest{
			Source:      source,
			Destination: dest,
		})
		if err != nil {
			return fmt.Errorf("failed to copy %q to %q: %w", source, dest, err)
		}

		fmt.Printf("Copied %s to %s\n", source, dest)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
