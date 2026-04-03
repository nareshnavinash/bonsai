package cmd

import (
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/chat"
)

var runCmd = &cobra.Command{
	Use:   "run [model] [prompt]",
	Short: "Start a chat session or run a one-shot prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		model := ""
		var promptParts []string

		if len(args) > 0 {
			// First arg: model name if it contains ":" or is a known short name
			first := args[0]
			if strings.Contains(first, ":") || len(args) > 1 {
				model = first
				promptParts = args[1:]
			} else if len(args) == 1 && len(first) < 30 && !strings.Contains(first, " ") {
				model = first
			} else {
				promptParts = args
			}
		}

		// If no model specified, pick the best available local model
		if model == "" {
			model, err = ResolveModel(client)
			if err != nil {
				return err
			}
		}

		prompt := strings.Join(promptParts, " ")

		systemMsg := api.Message{
			Role:    "system",
			Content: "You are a helpful assistant. Always respond in English unless the user explicitly asks for another language.",
		}

		if prompt != "" {
			// One-shot mode
			messages := []api.Message{
				systemMsg,
				{Role: "user", Content: prompt},
			}
			options := map[string]interface{}{
				"temperature": 0.7,
			}
			_, err := chat.StreamChat(client, model, messages, options)
			if err != nil {
				return err
			}
			fmt.Println()
			return nil
		}

		// Interactive REPL mode
		opts := &chat.REPLOptions{Temperature: 0.7}
		return chat.RunREPL(client, model, opts)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
