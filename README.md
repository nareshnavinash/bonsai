<div align="center">

<img src="assets/logo.svg" alt="bonsai" width="120">

# bonsai

**Run prism-ml's 1-bit Bonsai models locally via Ollama.**

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Ollama](https://img.shields.io/badge/Powered%20by-Ollama-black)](https://ollama.com)
[![Models by prism-ml](https://img.shields.io/badge/Models-prism--ml%20Bonsai-10b981)](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541)

[Quick Start](#quick-start) · [Models](#bonsai-models) · [Commands](#commands) · [Configuration](#configuration)

</div>

---

Bonsai is a CLI that makes it easy to run [prism-ml's Bonsai 1-bit quantized models](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541) through [Ollama](https://ollama.com). These models are 14x smaller than full-precision equivalents, use 4-5x less energy, and deliver fast inference on consumer hardware -- but Ollama doesn't natively surface them yet. Bonsai bridges that gap with a built-in model registry and a single `pull` command.

## Bonsai Models

The [prism-ml Bonsai models](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541) use **true 1-bit quantization** across all layers -- embeddings, attention, MLP, and output head. No escape hatches, no mixed-precision workarounds.

| Model | Parameters | Size | Pull Command |
|-------|-----------|------|-------------|
| **bonsai-8b** | 8B | ~1.2 GB | `bonsai pull bonsai-8b` |
| **bonsai-4b** | 4B | ~572 MB | `bonsai pull bonsai-4b` |
| **bonsai-1.7b** | 1.7B | ~248 MB | `bonsai pull bonsai-1.7b` |

### Why 1-Bit?

- **14x smaller** than FP16 equivalents
- **4-5x lower energy** consumption per token
- **Fast inference**: 40 tok/s on iPhone, 131 tok/s on M4 Pro, 368 tok/s on RTX 4090
- **Intelligence density**: 1.06 intelligence/GB vs 0.10 for full precision -- 10x more capability per byte
- **GGUF format** -- works natively with Ollama

Models by [prism-ml](https://huggingface.co/prism-ml) -- [explore the collection on HuggingFace](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541).

## Features

- **Built-in Bonsai registry** -- pull models by shortname (`bonsai pull bonsai-4b`), no need to remember HuggingFace paths
- **Full model management** -- pull, list, show, run, stop, remove, copy, create
- **Interactive chat** -- multi-turn conversations with streaming output
- **One-shot prompts** -- `bonsai run bonsai-4b "explain monads"`
- **Smart model resolution** -- auto-selects the best available Bonsai model
- **Thinking support** -- handles model reasoning transparently
- **Progress tracking** -- per-layer download progress bars
- **Lightweight** -- single binary, ~1,200 LOC, two dependencies

## Quick Start

### Prerequisites

[Ollama](https://ollama.com) must be installed and running:

```bash
curl -fsSL https://ollama.com/install.sh | sh
ollama serve
```

### Install Bonsai

```bash
go install github.com/nareshnavinash/bonsai@latest
```

Or download a binary from [Releases](https://github.com/nareshnavinash/bonsai/releases).

### Run

```bash
# Pull a model (~572 MB)
bonsai pull bonsai-4b

# Start chatting
bonsai run

# Or one-shot
bonsai run bonsai-4b "what is quantum computing?"
```

## Commands

| Command | Description |
|---------|-------------|
| `bonsai run [model] [prompt]` | Start a chat session or run a one-shot prompt |
| `bonsai pull <model>` | Download a model (supports bonsai shortnames) |
| `bonsai list` | List locally available models |
| `bonsai show <model>` | Show model details (arch, params, quantization) |
| `bonsai models` | List available Bonsai models from HuggingFace |
| `bonsai ps` | List currently running/loaded models |
| `bonsai stop <model>` | Unload a model from memory |
| `bonsai rm <model>` | Remove a model |
| `bonsai cp <src> <dest>` | Copy a model |
| `bonsai create <name> --from <model>` | Create a model with custom system prompt |
| `bonsai serve` | Start the Ollama server |
| `bonsai status` | Show server status and active model |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `BONSAI_MODEL` | `bonsai-8b` | Preferred model for commands |
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama server address |

### Model Resolution Order

1. `BONSAI_MODEL` environment variable (if set)
2. Any locally installed model with "bonsai" in its name
3. Any locally installed model
4. Helpful error message with pull instructions

### Chat Commands

In interactive mode (`bonsai run`):

| Command | Description |
|---------|-------------|
| `/bye` or `/exit` | Exit the chat |
| `/set temperature <value>` | Adjust creativity (0.0-2.0) |
| `/set top_p <value>` | Adjust nucleus sampling |
| `/set num_ctx <value>` | Adjust context window size |
| `"""` | Start multi-line input (end with `"""`) |

## Using Other Models

Bonsai works with any Ollama-compatible model:

```bash
bonsai pull llama3.2
bonsai run llama3.2 "explain monads"
BONSAI_MODEL=mistral bonsai run
```

## Contributing

Contributions are welcome. Please open an issue first to discuss what you would like to change.

```bash
git clone https://github.com/nareshnavinash/bonsai.git
cd bonsai
go build -o bonsai .
./bonsai status
```

## License

[MIT](LICENSE)

## Acknowledgments

- **[prism-ml](https://huggingface.co/prism-ml)** for the Bonsai 1-bit quantized model family
- **[Ollama](https://ollama.com)** for the local model runtime
- **[Cobra](https://github.com/spf13/cobra)** for the CLI framework
