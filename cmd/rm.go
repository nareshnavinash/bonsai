package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <model>",
	Short: "Remove a model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		model := args[0]

		client, err := getClient()
		if err != nil {
			return err
		}

		err = client.Delete(context.Background(), &api.DeleteRequest{Model: model})
		if err != nil {
			return fmt.Errorf("failed to delete %q: %w", model, err)
		}

		fmt.Printf("Deleted %s\n", model)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
