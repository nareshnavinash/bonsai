package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolErrorCmd = &cobra.Command{
	Use:   "error <message>",
	Short: "Explain an error or stack trace",
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
			return fmt.Errorf("usage: bonsai tool error <error message>")
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			`You are an error explainer. Given an error message or stack trace, explain:
1. What happened (one sentence)
2. Likely cause
3. Direction to fix (but do NOT write code)
Be concise (3-5 sentences).`,
			input)
	},
}

func init() {
	toolCmd.AddCommand(toolErrorCmd)
}
