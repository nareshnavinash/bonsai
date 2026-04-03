package main

import (
	"os"

	"github.com/nareshnavinash/bonsai/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
