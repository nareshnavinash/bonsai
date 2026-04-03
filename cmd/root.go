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
	Short: "Local model manager and coding assistant powered by Ollama",
	Long: `bonsai - Local model manager and coding assistant powered by Ollama

Model Commands:
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

Coding Tools:
  tool ask <question>          Ask a coding question
  tool summarize < file.py     Summarize code
  tool docstring < func.py     Generate a docstring
  tool error <message>         Explain an error
  tool pr                      Draft a PR description
  tool diff                    Explain staged diff
  tool regex <description>     Generate a regex
  tool commit                  Draft a commit message
  tool testdata <description>  Generate test data
  tool format <fmt> < input    Reformat text

Environment:
  BONSAI_MODEL      Model name (default: bonsai-8b)
  OLLAMA_HOST       Ollama server URL (default: http://localhost:11434)`,
	Version: "1.0.0",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

// ResolveModel picks the best model to use when none is specified.
// Prefers a locally installed bonsai model, falls back to any local model,
// or exits with a helpful message if nothing is installed.
func ResolveModel(client *api.Client) string {
	// If BONSAI_MODEL is set explicitly, use it
	if defaultModel != "bonsai-8b" {
		return defaultModel
	}

	resp, err := client.List(context.Background())
	if err != nil {
		return defaultModel
	}

	// First pass: look for any bonsai model
	for _, m := range resp.Models {
		if strings.Contains(strings.ToLower(m.Name), "bonsai") {
			return m.Name
		}
	}

	// Second pass: use any available model
	if len(resp.Models) > 0 {
		return resp.Models[0].Name
	}

	// Nothing installed
	fmt.Println("No models installed. Pull one first:")
	fmt.Println("  bonsai models        List available Bonsai models")
	fmt.Println("  bonsai pull bonsai-4b Download a model")
	os.Exit(1)
	return ""
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
