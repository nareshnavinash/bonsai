package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type BonsaiModel struct {
	Name     string
	HFRepo   string // e.g. "prism-ml/Bonsai-8B-gguf"
	GGUFFile string // e.g. "Bonsai-8B.gguf"
	Size     string
	Params   string
	Format   string
}

var Models = []BonsaiModel{
	{
		Name:     "bonsai-8b",
		HFRepo:   "prism-ml/Bonsai-8B-gguf",
		GGUFFile: "Bonsai-8B.gguf",
		Size:     "1.2 GB",
		Params:   "8B",
		Format:   "gguf",
	},
	{
		Name:     "bonsai-4b",
		HFRepo:   "prism-ml/Bonsai-4B-gguf",
		GGUFFile: "Bonsai-4B.gguf",
		Size:     "572 MB",
		Params:   "4B",
		Format:   "gguf",
	},
	{
		Name:     "bonsai-1.7b",
		HFRepo:   "prism-ml/Bonsai-1.7B-gguf",
		GGUFFile: "Bonsai-1.7B.gguf",
		Size:     "248 MB",
		Params:   "1.7B",
		Format:   "gguf",
	},
}

// FindModel looks up a model by name (case-insensitive, strips :tag).
// Returns nil if not a known bonsai model.
func FindModel(name string) *BonsaiModel {
	base := name
	if idx := strings.Index(name, ":"); idx != -1 {
		base = name[:idx]
	}
	for i := range Models {
		if strings.EqualFold(base, Models[i].Name) {
			return &Models[i]
		}
	}
	return nil
}

// Resolve maps a short bonsai name to its display name.
// Returns the original name unchanged if not a known bonsai model.
func Resolve(name string) string {
	m := FindModel(name)
	if m != nil {
		return m.Name
	}
	return name
}

// DownloadURL returns the direct HuggingFace download URL for this model.
func (m *BonsaiModel) DownloadURL() string {
	return fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", m.HFRepo, m.GGUFFile)
}

// LocalPath returns the expected local file path for this model.
func (m *BonsaiModel) LocalPath() string {
	return filepath.Join(ModelsDir(), m.GGUFFile)
}

// ModelsDir returns the models directory path.
// Uses BONSAI_MODELS_DIR env var, defaulting to ~/.bonsai/models/
func ModelsDir() string {
	if dir := os.Getenv("BONSAI_MODELS_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bonsai", "models")
}
