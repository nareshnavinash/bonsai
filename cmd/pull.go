package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/ui"
)

var pullCmd = &cobra.Command{
	Use:   "pull <model>",
	Short: "Download a model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		model := registry.FindModel(args[0])
		if model == nil {
			return fmt.Errorf("unknown model %q. Run 'bonsai models' to see available models", args[0])
		}

		// Check if already downloaded
		if registry.IsDownloaded(model) {
			fmt.Printf("%s is already downloaded at %s\n", model.Name, model.LocalPath())
			return nil
		}

		fmt.Printf("Downloading %s (%s)...\n", model.Name, model.Size)

		var bar *ui.ProgressBar

		err := registry.Download(cmd.Context(), model, func(downloaded, total int64) {
			if bar == nil && total > 0 {
				bar = ui.NewProgressBar(total, "pulling... ")
			}
			if bar != nil {
				bar.Update(downloaded)
			}
		})
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		if bar != nil {
			bar.Done()
		}
		fmt.Printf("Downloaded %s to %s\n", model.Name, model.LocalPath())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
