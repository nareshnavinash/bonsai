package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolSummarizeCmd = &cobra.Command{
	Use:   "summarize [text]",
	Short: "Summarize code (reads from stdin if piped)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else if tool.IsPiped() {
			var err error
			input, err = tool.ReadStdin()
			if err != nil {
				return err
			}
		}
		if input == "" {
			return fmt.Errorf("no input provided. Pipe code or provide text as arguments")
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			"You are a code explainer. Summarize what the code does in 2-5 sentences. Focus on WHAT, not HOW.",
			input)
	},
}

func init() {
	toolCmd.AddCommand(toolSummarizeCmd)
}
