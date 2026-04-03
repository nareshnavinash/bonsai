package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <model>",
	Short: "Unload a model from memory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		model := args[0]

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		keepAlive := &api.Duration{Duration: 0}
		err = client.Generate(context.Background(), &api.GenerateRequest{
			Model:     model,
			KeepAlive: keepAlive,
		}, func(resp api.GenerateResponse) error {
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to stop %q: %w", model, err)
		}

		fmt.Printf("Stopped %s\n", model)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
