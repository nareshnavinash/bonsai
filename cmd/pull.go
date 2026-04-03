package cmd

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
	"github.com/nareshnavinash/bonsai/internal/ui"
)

var pullCmd = &cobra.Command{
	Use:   "pull <model>",
	Short: "Download a model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		model := registry.Resolve(args[0])

		client, err := getClient()
		if err != nil {
			return err
		}

		var currentBar *ui.ProgressBar
		var currentDigest string

		defer func() {
			if currentBar != nil {
				currentBar.Done()
			}
		}()

		err = client.Pull(context.Background(), &api.PullRequest{Model: model}, func(resp api.ProgressResponse) error {
			if resp.Digest != "" && resp.Total > 0 {
				if resp.Digest != currentDigest {
					if currentBar != nil {
						currentBar.Done()
					}
					shortDigest := ui.TruncateID(resp.Digest)
					currentBar = ui.NewProgressBar(resp.Total, fmt.Sprintf("pulling %s... ", shortDigest))
					currentDigest = resp.Digest
				}
				if currentBar != nil {
					currentBar.Update(resp.Completed)
				}
			} else {
				if currentBar != nil {
					currentBar.Done()
					currentBar = nil
					currentDigest = ""
				}
				fmt.Println(resp.Status)
			}
			return nil
		})

		return err
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
