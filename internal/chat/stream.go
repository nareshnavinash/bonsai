package chat

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nareshnavinash/bonsai/internal/llm"
)

func filterArtifacts(s string) string {
	s = strings.ReplaceAll(s, "<tool_call>", "")
	s = strings.ReplaceAll(s, "</tool_call>", "")
	return strings.TrimSpace(s)
}

func StreamChat(client *llm.Client, model string, messages []llm.Message, options map[string]interface{}) (string, error) {
	req := &llm.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	// Apply options
	if v, ok := options["temperature"]; ok {
		if f, ok := v.(float64); ok {
			req.Temperature = &f
		}
	}
	if v, ok := options["top_p"]; ok {
		if f, ok := v.(float64); ok {
			req.TopP = &f
		}
	}

	var response string

	err := client.ChatCompletionStream(context.Background(), req, func(resp llm.ChatCompletionResponse) error {
		if len(resp.Choices) > 0 && resp.Choices[0].Delta != nil {
			content := resp.Choices[0].Delta.Content
			if content != "" {
				fmt.Fprint(os.Stdout, content)
				response += content
			}
		}
		return nil
	})

	// Filter artifacts from accumulated response
	response = filterArtifacts(response)

	return response, err
}
