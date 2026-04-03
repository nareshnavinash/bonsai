package registry

import "strings"

type BonsaiModel struct {
	Name     string
	HFRepo   string
	Size     string
	Params   string
	Format   string
}

var Models = []BonsaiModel{
	{
		Name:   "bonsai-8b",
		HFRepo: "hf.co/prism-ml/Bonsai-8B-gguf",
		Size:   "1.2 GB",
		Params: "8B",
		Format: "gguf",
	},
	{
		Name:   "bonsai-4b",
		HFRepo: "hf.co/prism-ml/Bonsai-4B-gguf",
		Size:   "572 MB",
		Params: "4B",
		Format: "gguf",
	},
	{
		Name:   "bonsai-1.7b",
		HFRepo: "hf.co/prism-ml/Bonsai-1.7B-gguf",
		Size:   "248 MB",
		Params: "1.7B",
		Format: "gguf",
	},
}

// Resolve maps a short bonsai name (e.g. "bonsai-4b") to its HuggingFace
// pull path. Returns the original name unchanged if not a known bonsai model.
func Resolve(name string) string {
	// Strip tag if present for matching (e.g. "bonsai-4b:latest" -> "bonsai-4b")
	base := name
	if idx := strings.Index(name, ":"); idx != -1 {
		base = name[:idx]
	}

	for _, m := range Models {
		if strings.EqualFold(base, m.Name) {
			return m.HFRepo
		}
	}
	return name
}
