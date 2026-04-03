package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/tool"
)

var toolFormatCmd = &cobra.Command{
	Use:   "format <target-format> [text]",
	Short: "Reformat text to a target format",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: bonsai tool format <target-format> < input.txt")
		}

		targetFormat := args[0]
		var text string
		if len(args) > 1 {
			text = strings.Join(args[1:], " ")
		} else if tool.IsPiped() {
			var err error
			text, err = tool.ReadStdin()
			if err != nil {
				return err
			}
		}
		if text == "" {
			return fmt.Errorf("no input provided. Pipe data or provide text as arguments")
		}

		client, err := api.ClientFromEnvironment()
		if err != nil {
			return err
		}

		return tool.StreamWithSystemPrompt(client, ResolveModel(client),
			"You are a text formatter. Convert the input to the requested format. Output ONLY the result, no explanations.",
			fmt.Sprintf("Convert to %s:\n\n%s", targetFormat, text))
	},
}

func init() {
	toolCmd.AddCommand(toolFormatCmd)
}
