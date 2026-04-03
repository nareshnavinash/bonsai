package chat

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ollama/ollama/api"
)

type REPLOptions struct {
	Temperature float64
	TopP        float64
	NumCtx      int
}

func RunREPL(client *api.Client, model string, opts *REPLOptions) error {
	scanner := bufio.NewScanner(os.Stdin)
	messages := []api.Message{
		{Role: "system", Content: "You are a helpful assistant. Always respond in English unless the user explicitly asks for another language."},
	}

	fmt.Printf("Interactive chat with %s. Type /bye to exit, \"\"\" for multi-line input.\n\n", model)

	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "/bye" || input == "/exit" {
			break
		}

		// Handle /set command
		if strings.HasPrefix(input, "/set ") {
			handleSet(input[5:], opts)
			continue
		}

		// Handle multi-line input
		fullInput := input
		if strings.HasPrefix(input, `"""`) {
			fullInput = input[3:]
			if strings.HasSuffix(fullInput, `"""`) {
				fullInput = fullInput[:len(fullInput)-3]
			} else {
				fullInput += "\n"
				for {
					fmt.Print("... ")
					if !scanner.Scan() {
						break
					}
					line := scanner.Text()
					if strings.HasSuffix(line, `"""`) {
						fullInput += line[:len(line)-3]
						break
					}
					fullInput += line + "\n"
				}
			}
		}

		messages = append(messages, api.Message{Role: "user", Content: fullInput})

		options := buildOptions(opts)
		response, err := StreamChat(client, model, messages, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
			// Remove failed user message
			messages = messages[:len(messages)-1]
			continue
		}

		fmt.Print("\n\n")
		if strings.TrimSpace(response) != "" {
			messages = append(messages, api.Message{Role: "assistant", Content: response})
		}
	}

	return nil
}

func handleSet(input string, opts *REPLOptions) {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		fmt.Println("Usage: /set <temperature|top_p|num_ctx> <value>")
		return
	}
	key, val := parts[0], parts[1]
	switch key {
	case "temperature":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			opts.Temperature = f
			fmt.Printf("Set temperature to %.1f\n", f)
		}
	case "top_p":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			opts.TopP = f
			fmt.Printf("Set top_p to %.1f\n", f)
		}
	case "num_ctx":
		if n, err := strconv.Atoi(val); err == nil {
			opts.NumCtx = n
			fmt.Printf("Set num_ctx to %d\n", n)
		}
	default:
		fmt.Println("Usage: /set <temperature|top_p|num_ctx> <value>")
	}
}

func buildOptions(opts *REPLOptions) map[string]interface{} {
	m := map[string]interface{}{}
	if opts.Temperature > 0 {
		m["temperature"] = opts.Temperature
	}
	if opts.TopP > 0 {
		m["top_p"] = opts.TopP
	}
	if opts.NumCtx > 0 {
		m["num_ctx"] = opts.NumCtx
	}
	return m
}
