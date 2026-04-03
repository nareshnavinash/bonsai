package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nareshnavinash/bonsai/internal/registry"
)

var cpCmd = &cobra.Command{
	Use:   "cp <source> <destination>",
	Short: "Copy a model file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath, err := registry.ResolveModelPath(args[0])
		if err != nil {
			return err
		}

		// If destination is a short name, put it in models dir
		dstPath := args[1]
		if !strings.Contains(dstPath, "/") && !strings.HasSuffix(dstPath, ".gguf") {
			dstPath = filepath.Join(registry.ModelsDir(), dstPath+".gguf")
		}

		src, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("cannot open source: %w", err)
		}
		defer src.Close()

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("cannot create destination: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("copy failed: %w", err)
		}

		fmt.Printf("Copied %s to %s\n", srcPath, dstPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
