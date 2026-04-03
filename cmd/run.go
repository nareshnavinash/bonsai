package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/chat"
	"github.com/nareshnavinash/bonsai/internal/llm"
	"github.com/nareshnavinash/bonsai/internal/registry"
)

var runCmd = &cobra.Command{
	Use:   "run [model] [prompt]",
	Short: "Start a chat session or run a one-shot prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		model := ""
		var promptParts []string

		if len(args) > 0 {
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

		if model == "" {
			var err error
			model, err = ResolveModel()
			if err != nil {
				return err
			}
		}

		// Resolve model name to GGUF file path
		modelPath, err := registry.ResolveModelPath(model)
		if err != nil {
			return err
		}

		// Ensure llama-server is running with this model
		mgr := getServerManager()
		baseURL, err := mgr.EnsureRunning(modelPath)
		if err != nil {
			return err
		}

		client := llm.NewClient(baseURL)

		prompt := strings.Join(promptParts, " ")

		systemMsg := llm.Message{
			Role:    "system",
			Content: "You are a helpful assistant. Always respond in English unless the user explicitly asks for another language.",
		}

		if prompt != "" {
			// One-shot mode
			messages := []llm.Message{
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
		switchFn := func(name string) (string, error) {
			newPath, err := registry.ResolveModelPath(name)
			if err != nil {
				return "", err
			}
			_, err = mgr.EnsureRunning(newPath)
			if err != nil {
				return "", err
			}
			return registry.Resolve(name), nil
		}
		return chat.RunREPL(client, model, opts, switchFn)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
