package cmd

import (
	"github.com/spf13/cobra"
)

var toolCmd = &cobra.Command{
	Use:   "tool <subcommand>",
	Short: "Coding tools (ask, summarize, docstring, error, pr, diff, regex, commit, testdata, format)",
}

func init() {
	rootCmd.AddCommand(toolCmd)
}
