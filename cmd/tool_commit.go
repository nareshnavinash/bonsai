package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolCommitCmd = &cobra.Command{
	Use:   "commit [text]",
	Short: "Draft a commit message from staged changes",
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
			return fmt.Errorf("no staged changes. Stage files with 'git add' first")
		}

		if len(input) > 3000 {
			input = input[:3000] + "\n...(truncated)"
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			`You are a commit message writer. Write a conventional commit message from a git diff.
Format: type(scope): description. Types: feat, fix, refactor, test, docs, chore.
Subject under 72 chars. Output ONLY the message.`,
			input)
	},
}

func init() {
	toolCmd.AddCommand(toolCommitCmd)
}
