package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var showModelfile bool
var showParameters bool
var showSystem bool

var showCmd = &cobra.Command{
	Use:   "show <model>",
	Short: "Show model details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		model := args[0]

		client, err := getClient()
		if err != nil {
			return err
		}

		resp, err := client.Show(context.Background(), &api.ShowRequest{Model: model})
		if err != nil {
			return fmt.Errorf("cannot show model %q: %w", model, err)
		}

		if showModelfile {
			fmt.Println(resp.Modelfile)
			return nil
		}

		if showParameters {
			if resp.Parameters != "" {
				fmt.Println(resp.Parameters)
			} else {
				fmt.Println("No parameters set.")
			}
			return nil
		}

		if showSystem {
			if resp.System != "" {
				fmt.Println(resp.System)
			} else {
				fmt.Println("No system prompt set.")
			}
			return nil
		}

		// Default: structured overview
		fmt.Println("Model")
		if resp.Details.Family != "" {
			fmt.Printf("  arch            %s\n", resp.Details.Family)
		}
		if resp.Details.ParameterSize != "" {
			fmt.Printf("  parameters      %s\n", resp.Details.ParameterSize)
		}
		if resp.Details.QuantizationLevel != "" {
			fmt.Printf("  quantization    %s\n", resp.Details.QuantizationLevel)
		}
		if resp.Details.Format != "" {
			fmt.Printf("  format          %s\n", resp.Details.Format)
		}
		fmt.Println()

		if resp.Parameters != "" {
			fmt.Println("Parameters")
			for _, line := range strings.Split(resp.Parameters, "\n") {
				if strings.TrimSpace(line) != "" {
					fmt.Printf("  %s\n", strings.TrimSpace(line))
				}
			}
			fmt.Println()
		}

		if resp.System != "" {
			fmt.Println("System")
			fmt.Printf("  %s\n", resp.System)
		}

		return nil
	},
}

func init() {
	showCmd.Flags().BoolVar(&showModelfile, "modelfile", false, "Show raw Modelfile")
	showCmd.Flags().BoolVar(&showParameters, "parameters", false, "Show parameters only")
	showCmd.Flags().BoolVar(&showSystem, "system", false, "Show system prompt only")
	rootCmd.AddCommand(showCmd)
}
