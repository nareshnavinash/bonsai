package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolAskCmd = &cobra.Command{
	Use:   "ask <question>",
	Short: "Ask a coding question",
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")
		if question == "" {
			return fmt.Errorf("usage: bonsai tool ask <question>")
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			"You are a helpful coding assistant. Be concise and direct.",
			question)
	},
}

func init() {
	toolCmd.AddCommand(toolAskCmd)
}
