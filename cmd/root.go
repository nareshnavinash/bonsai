package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var defaultModel = "bonsai-8b"

func init() {
	if m := os.Getenv("BONSAI_MODEL"); m != "" {
		defaultModel = m
	}
}

var rootCmd = &cobra.Command{
	Use:   "bonsai",
	Short: "Run prism-ml's 1-bit Bonsai models locally via Ollama",
	Long: `bonsai - Run prism-ml's 1-bit Bonsai models locally via Ollama

Commands:
  run [model] [prompt]     Start a chat session or run a one-shot prompt
  pull <model>             Download a model
  list                     List available models
  show <model>             Show model details
  ps                       List running models
  stop <model>             Unload a model from memory
  rm <model>               Remove a model
  cp <source> <dest>       Copy a model
  create <name> -f <file>  Create a model from a Modelfile
  serve                    Start the Ollama server
  models                   List available Bonsai models
  status                   Show server status

Environment:
  BONSAI_MODEL      Model name (default: bonsai-8b)
  OLLAMA_HOST       Ollama server URL (default: http://localhost:11434)`,
	Version: "1.0.0",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

// getClient creates an Ollama API client with a helpful error message on failure.
func getClient() (*api.Client, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Ollama (is it running? try 'bonsai serve'): %w", err)
	}
	return client, nil
}

// ResolveModel picks the best model to use when none is specified.
// Prefers a locally installed bonsai model, falls back to any local model.
func ResolveModel(client *api.Client) (string, error) {
	// If BONSAI_MODEL is set explicitly, use it
	if defaultModel != "bonsai-8b" {
		return defaultModel, nil
	}

	resp, err := client.List(context.Background())
	if err != nil {
		return defaultModel, nil
	}

	// First pass: look for any bonsai model
	for _, m := range resp.Models {
		if strings.Contains(strings.ToLower(m.Name), "bonsai") {
			return m.Name, nil
		}
	}

	// Second pass: use any available model
	if len(resp.Models) > 0 {
		return resp.Models[0].Name, nil
	}

	return "", fmt.Errorf("no models installed. Pull one first:\n  bonsai models        List available Bonsai models\n  bonsai pull bonsai-4b Download a model")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
