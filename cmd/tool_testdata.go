package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolTestdataCmd = &cobra.Command{
	Use:   "testdata <description>",
	Short: "Generate test data from a description",
	RunE: func(cmd *cobra.Command, args []string) error {
		desc := strings.Join(args, " ")
		if desc == "" {
			return fmt.Errorf("usage: bonsai tool testdata <description>")
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			"You are a test data generator. Generate realistic but fictional test data. Use obviously fake values. Output valid JSON.",
			desc)
	},
}

func init() {
	toolCmd.AddCommand(toolTestdataCmd)
}
