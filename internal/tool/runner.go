package tool

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/nareshnavinash/bonsai/internal/chat"
)

func StreamWithSystemPrompt(client *api.Client, model, systemPrompt, userInput string) error {
	messages := []api.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userInput},
	}

	options := map[string]interface{}{
		"temperature": 0.3,
	}

	_, err := chat.StreamChat(client, model, messages, options)
	if err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func ReadStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func IsPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
