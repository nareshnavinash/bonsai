package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshsekar/bonsai/internal/tool"
)

var toolRegexCmd = &cobra.Command{
	Use:   "regex <description>",
	Short: "Generate a regex from a description",
	RunE: func(cmd *cobra.Command, args []string) error {
		desc := strings.Join(args, " ")
		if desc == "" {
			return fmt.Errorf("usage: bonsai tool regex <description>")
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			`You are a regex generator. Given a description, output:
Pattern: <regex>
Explanation: <brief explanation of each part>
Example matches: <2-3 examples>`,
			desc)
	},
}

func init() {
	toolCmd.AddCommand(toolRegexCmd)
}
