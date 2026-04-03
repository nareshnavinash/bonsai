package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalModel represents a GGUF file found on disk.
type LocalModel struct {
	Name       string
	Path       string
	Size       int64
	ModifiedAt time.Time
	Known      bool // true if it matches a registry entry
}

// ScanLocal scans the models directory for .gguf files.
func ScanLocal() ([]LocalModel, error) {
	dir := ModelsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var models []LocalModel
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".gguf") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}

		name := e.Name()
		known := false
		for _, m := range Models {
			if strings.EqualFold(m.GGUFFile, e.Name()) {
				name = m.Name
				known = true
				break
			}
		}

		models = append(models, LocalModel{
			Name:       name,
			Path:       filepath.Join(dir, e.Name()),
			Size:       info.Size(),
			ModifiedAt: info.ModTime(),
			Known:      known,
		})
	}

	return models, nil
}

// ResolveModelPath takes a model name (e.g. "bonsai-8b") or a file path
// and returns the absolute path to the .gguf file.
func ResolveModelPath(name string) (string, error) {
	// 1. Absolute or relative path that exists
	if filepath.IsAbs(name) || strings.Contains(name, "/") || strings.Contains(name, ".gguf") {
		if _, err := os.Stat(name); err == nil {
			abs, _ := filepath.Abs(name)
			return abs, nil
		}
	}

	// 2. Registry lookup — check standard local path
	if m := FindModel(name); m != nil {
		path := m.LocalPath()
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 3. Backward compat: check ~/models/bonsai-*/ directories
	home, _ := os.UserHomeDir()
	legacyPatterns := []string{
		filepath.Join(home, "models", name, "*.gguf"),
		filepath.Join(home, "models", name+"*", "*.gguf"),
	}
	for _, pattern := range legacyPatterns {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			return matches[0], nil
		}
	}

	// 4. Scan models dir for fuzzy filename match
	dir := ModelsDir()
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Name()), strings.ToLower(name)) {
			return filepath.Join(dir, e.Name()), nil
		}
	}

	return "", fmt.Errorf("model %q not found locally.\n  Pull it first: bonsai pull %s\n  Or specify a path to a .gguf file", name, name)
}
