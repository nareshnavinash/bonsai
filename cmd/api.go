package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	bonsaiapi "github.com/nareshnavinash/bonsai/internal/api"
)

var (
	apiPort int
	apiHost string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start an OpenAI-compatible API server",
	Long: `Start an HTTP server exposing OpenAI-compatible endpoints.

Requires Ollama to be running (bonsai serve or ollama serve).

Endpoints:
  POST /v1/chat/completions   Chat completions (streaming & non-streaming)
  GET  /v1/models             List available models
  GET  /health                Health check

Usage with OpenAI SDK:
  from openai import OpenAI
  client = OpenAI(base_url="http://localhost:8080/v1", api_key="unused")`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return err
		}

		model, err := ResolveModel(client)
		if err != nil {
			return err
		}

		fmt.Printf("Default model: %s\n", model)

		server := &bonsaiapi.Server{
			Client:       client,
			DefaultModel: model,
			Host:         apiHost,
			Port:         apiPort,
		}

		return server.Start()
	},
}

func init() {
	apiCmd.Flags().IntVar(&apiPort, "port", 8080, "Port to listen on")
	apiCmd.Flags().StringVar(&apiHost, "host", "localhost", "Host to bind to")
	rootCmd.AddCommand(apiCmd)
}
