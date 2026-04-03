package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/llm"
	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/server"
)

var defaultModel = "bonsai-8b"

func init() {
	if m := os.Getenv("BONSAI_MODEL"); m != "" {
		defaultModel = m
	}
}

var rootCmd = &cobra.Command{
	Use:   "bonsai",
	Short: "Run prism-ml's 1-bit Bonsai models locally via llama.cpp",
	Long: `bonsai - Run prism-ml's 1-bit Bonsai models locally

Commands:
  run [model] [prompt]     Start a chat session or run a one-shot prompt
  pull <model>             Download a model
  list                     List installed models
  show <model>             Show model details
  ps                       Show running server status
  stop                     Stop the server
  rm <model>               Remove a model
  cp <source> <dest>       Copy a model file
  serve [model]            Start the llama-server
  api                      Start OpenAI-compatible API server
  models                   List available Bonsai models
  status                   Show server status

Environment:
  BONSAI_MODEL        Model name (default: bonsai-8b)
  BONSAI_HOST         Server URL (default: http://127.0.0.1:8081)
  BONSAI_PORT         Server port (default: 8081)
  BONSAI_THREADS      CPU threads for inference
  BONSAI_MODELS_DIR   Model storage directory (default: ~/.bonsai/models/)
  LLAMA_SERVER_BIN    Path to llama-server binary`,
	Version: "2.0.0",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

// getLLMClient creates an OpenAI-compatible HTTP client.
func getLLMClient() *llm.Client {
	return llm.NewClientFromEnv()
}

// getServerManager creates a server lifecycle manager.
func getServerManager() *server.Manager {
	return server.NewManager()
}

// ResolveModel picks the best model to use when none is specified.
func ResolveModel() (string, error) {
	if defaultModel != "bonsai-8b" {
		return defaultModel, nil
	}

	locals, err := registry.ScanLocal()
	if err != nil {
		return defaultModel, nil
	}

	// First pass: look for any known bonsai model
	for _, m := range locals {
		if m.Known {
			return m.Name, nil
		}
	}

	// Second pass: use any available .gguf file
	if len(locals) > 0 {
		return locals[0].Name, nil
	}

	// Check legacy paths
	for _, m := range registry.Models {
		if _, err := registry.ResolveModelPath(m.Name); err == nil {
			return m.Name, nil
		}
	}

	return "", fmt.Errorf("no models installed. Pull one first:\n  bonsai models        List available Bonsai models\n  bonsai pull bonsai-4b Download a model")
}

// ResolveModelForDisplay returns the model name with any name resolution applied.
func ResolveModelForDisplay(name string) string {
	return strings.TrimSuffix(registry.Resolve(name), ":latest")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
