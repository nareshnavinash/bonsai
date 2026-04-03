package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/ui"
)

var showCmd = &cobra.Command{
	Use:   "show <model>",
	Short: "Show model details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Check registry
		model := registry.FindModel(name)
		if model != nil {
			fmt.Println("Model")
			fmt.Printf("  name            %s\n", model.Name)
			fmt.Printf("  parameters      %s\n", model.Params)
			fmt.Printf("  format          %s\n", model.Format)
			fmt.Printf("  registry size   %s\n", model.Size)
			fmt.Printf("  huggingface     https://huggingface.co/%s\n", model.HFRepo)
			fmt.Println()
		}

		// Check local file
		localPath := ""
		if model != nil {
			localPath = model.LocalPath()
		} else {
			localPath, _ = registry.ResolveModelPath(name)
		}

		if localPath != "" {
			if info, err := os.Stat(localPath); err == nil {
				fmt.Println("Local")
				fmt.Printf("  path            %s\n", localPath)
				fmt.Printf("  file size       %s\n", ui.FormatBytes(info.Size()))
				fmt.Printf("  modified        %s\n", ui.FormatRelativeTime(info.ModTime()))
			} else if model != nil {
				fmt.Println("Status: not downloaded")
				fmt.Printf("  pull with:      bonsai pull %s\n", name)
			}
		}

		if model == nil && localPath == "" {
			return fmt.Errorf("model %q not found", name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
