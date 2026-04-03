package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolPRCmd = &cobra.Command{
	Use:   "pr [text]",
	Short: "Draft a PR description from a diff or commit log",
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else {
			data, err := tool.CommitLog()
			if err != nil {
				return fmt.Errorf("not a git repository or git is not available")
			}
			input = data
		}
		if input == "" {
			return fmt.Errorf("no changes found")
		}

		if len(input) > 3000 {
			input = input[:3000] + "\n...(truncated)"
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			`You are a PR description writer. Given a diff or commit log, output:
## Title
<short title under 70 chars>
## Summary
<2-4 bullet points>
## Test plan
<2-3 bullet points>`,
			input)
	},
}

func init() {
	toolCmd.AddCommand(toolPRCmd)
}
