package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// FindBinary locates the llama-server binary.
// Search order:
//  1. LLAMA_SERVER_BIN env var
//  2. "llama-server" in PATH
//  3. Common installation paths
func FindBinary() (string, error) {
	// 1. Environment variable
	if bin := os.Getenv("LLAMA_SERVER_BIN"); bin != "" {
		if _, err := os.Stat(bin); err == nil {
			return bin, nil
		}
		return "", fmt.Errorf("LLAMA_SERVER_BIN set to %q but file not found", bin)
	}

	// 2. PATH lookup
	if path, err := exec.LookPath("llama-server"); err == nil {
		return path, nil
	}

	// 3. Common locations
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, "projects", "llama.cpp", "build", "bin", "llama-server"),
		"/usr/local/bin/llama-server",
		"/opt/homebrew/bin/llama-server",
		filepath.Join(home, ".local", "bin", "llama-server"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("llama-server not found. Install it:\n" +
		"  brew install llama.cpp\n" +
		"  or: cd ~/projects/llama.cpp && cmake -B build && cmake --build build -j\n" +
		"  or: set LLAMA_SERVER_BIN to the binary path")
}
