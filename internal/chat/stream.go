package chat

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/nareshsekar/bonsai/internal/ui"
)

func filterArtifacts(s string) string {
	// Remove common model artifacts
	s = strings.ReplaceAll(s, "<tool_call>", "")
	s = strings.ReplaceAll(s, "</tool_call>", "")
	return strings.TrimSpace(s)
}

func StreamChat(client *api.Client, model string, messages []api.Message, options map[string]interface{}) (string, error) {
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Options:  options,
	}

	var thinkBuf string
	var response string
	var gotContent bool
	var spinner *ui.Spinner

	err := client.Chat(context.Background(), req, func(resp api.ChatResponse) error {
		if resp.Message.Content != "" {
			if !gotContent {
				// First content token — stop spinner, start streaming
				gotContent = true
				if spinner != nil {
					spinner.Stop("")
					spinner = nil
				}
			}
			fmt.Fprint(os.Stdout, resp.Message.Content)
			response += resp.Message.Content
		} else if resp.Message.Thinking != "" {
			// Buffer thinking silently, show spinner
			if spinner == nil && !gotContent {
				spinner = ui.NewSpinner("thinking...")
				spinner.Start()
			}
			thinkBuf += resp.Message.Thinking
		}
		return nil
	})

	if spinner != nil {
		spinner.Stop("")
	}

	// If no content tokens arrived, the thinking IS the response
	if !gotContent && thinkBuf != "" {
		// Filter out model artifacts like <tool_call> tags
		cleaned := filterArtifacts(thinkBuf)
		if cleaned != "" {
			fmt.Fprint(os.Stdout, cleaned)
			response = cleaned
		}
	}

	return response, err
}
