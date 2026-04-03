package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var (
	createFrom   string
	createSystem string
)

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a model from a base model",
	Long:  "Create a new model. Use --from to specify the base model and --system for a custom system prompt.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if createFrom == "" {
			return fmt.Errorf("base model required: bonsai create %s --from <model>", name)
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		req := &api.CreateRequest{
			Model: name,
			From:  createFrom,
		}
		if createSystem != "" {
			req.System = createSystem
		}

		err = client.Create(context.Background(), req, func(resp api.ProgressResponse) error {
			fmt.Println(resp.Status)
			return nil
		})

		return err
	},
}

func init() {
	createCmd.Flags().StringVar(&createFrom, "from", "", "Base model to create from")
	createCmd.Flags().StringVar(&createSystem, "system", "", "System prompt for the model")
	rootCmd.AddCommand(createCmd)
}
