package main

import (
	"os"

	"github.com/nareshsekar/bonsai/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
