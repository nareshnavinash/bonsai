package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshsekar/bonsai/internal/tool"
)

var toolDiffCmd = &cobra.Command{
	Use:   "diff [text]",
	Short: "Explain staged diff in plain English",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else {
			data, err := tool.StagedDiff()
			if err != nil {
				return fmt.Errorf("not a git repository or git is not available")
			}
			input = data
		}
		if input == "" {
			return fmt.Errorf("no staged changes found")
		}

		if len(input) > 3000 {
			input = input[:3000] + "\n...(truncated)"
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			`You are a diff explainer. Summarize what changed in 2-5 sentences.
Group related changes, mention files affected. Focus on WHAT and WHY. Do NOT write code.`,
			input)
	},
}

func init() {
	toolCmd.AddCommand(toolDiffCmd)
}
