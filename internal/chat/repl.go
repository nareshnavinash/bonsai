package chat

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nareshnavinash/bonsai/internal/llm"
	"github.com/nareshnavinash/bonsai/internal/registry"
)

type REPLOptions struct {
	Temperature float64
	TopP        float64
	NumCtx      int
}

// SwitchModelFn is called when the user runs /model <name> in the REPL.
// It should restart the server with the new model and return the new model name.
type SwitchModelFn func(name string) (newModel string, err error)

func RunREPL(client *llm.Client, model string, opts *REPLOptions, switchFn SwitchModelFn) error {
	scanner := bufio.NewScanner(os.Stdin)
	systemMsg := llm.Message{Role: "system", Content: "You are a helpful assistant. Always respond in English unless the user explicitly asks for another language."}
	messages := []llm.Message{systemMsg}

	fmt.Printf("Interactive chat with %s. Type /help for commands.\n\n", model)

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

		if input == "/help" {
			printHelp()
			continue
		}

		if input == "/clear" {
			messages = []llm.Message{systemMsg}
			fmt.Println("Conversation cleared.")
			continue
		}

		if strings.HasPrefix(input, "/model ") {
			newModel := strings.TrimSpace(input[7:])
			if newModel == "" {
				fmt.Fprintf(os.Stderr, "Usage: /model <name>\n")
				continue
			}
			if switchFn != nil {
				resolved, err := switchFn(newModel)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error switching model: %v\n", err)
					continue
				}
				model = resolved
			} else {
				model = registry.Resolve(newModel)
			}
			messages = []llm.Message{systemMsg}
			fmt.Printf("Switched to %s. Conversation cleared.\n", model)
			continue
		}

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

		messages = append(messages, llm.Message{Role: "user", Content: fullInput})

		options := buildOptions(opts)
		response, err := StreamChat(client, model, messages, options)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
			messages = messages[:len(messages)-1]
			continue
		}

		fmt.Print("\n\n")
		if strings.TrimSpace(response) != "" {
			messages = append(messages, llm.Message{Role: "assistant", Content: response})
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input error: %w", err)
	}

	return nil
}

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  /bye, /exit           Exit the chat")
	fmt.Println("  /clear                Clear conversation history")
	fmt.Println("  /model <name>         Switch to a different model")
	fmt.Println("  /set temperature <n>  Set temperature (0.0-2.0)")
	fmt.Println("  /set top_p <n>        Set top_p (0.0-1.0)")
	fmt.Println("  /set num_ctx <n>      Set context window size")
	fmt.Println(`  """                   Start multi-line input (end with """)`)
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
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid number %q for temperature\n", val)
			return
		}
		opts.Temperature = f
		fmt.Printf("Set temperature to %.1f\n", f)
	case "top_p":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid number %q for top_p\n", val)
			return
		}
		opts.TopP = f
		fmt.Printf("Set top_p to %.1f\n", f)
	case "num_ctx":
		n, err := strconv.Atoi(val)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid integer %q for num_ctx\n", val)
			return
		}
		opts.NumCtx = n
		fmt.Printf("Set num_ctx to %d\n", n)
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
	return m
}
