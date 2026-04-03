package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshsekar/bonsai/internal/tool"
)

var toolDocstringCmd = &cobra.Command{
	Use:   "docstring [code]",
	Short: "Generate a docstring (reads code from stdin)",
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
			"You are a docstring generator. Given Python code, produce a Google-style docstring. Output ONLY the docstring with triple quotes.",
			input)
	},
}

func init() {
	toolCmd.AddCommand(toolDocstringCmd)
}
