<div align="center">

<img src="assets/logo.svg" alt="bonsai" width="120">

# bonsai

**Local AI tools for developers. No cloud. No API keys. Just code.**

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Ollama](https://img.shields.io/badge/Powered%20by-Ollama-black)](https://ollama.com)
[![Models by prism-ml](https://img.shields.io/badge/Models-prism--ml%20Bonsai-10b981)](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541)

[Features](#features) · [Quick Start](#quick-start) · [Commands](#commands) · [Tools](#developer-tools) · [Models](#bonsai-models) · [Configuration](#configuration)

</div>

---

Bonsai wraps [Ollama](https://ollama.com) into a single binary with 10 developer-focused AI tools that run entirely on your machine. Draft commit messages from real diffs. Generate PR descriptions from your branch history. Explain errors. Summarize code. All locally, all private.

Ships with first-class support for the [prism-ml Bonsai models](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541) -- true 1-bit quantized language models that are 14x smaller than full-precision equivalents while maintaining strong coding performance.

## Features

- **10 developer tools** -- ask, summarize, docstring, error, commit, diff, pr, regex, testdata, format
- **Git-aware** -- reads staged diffs and branch history automatically
- **12 model management commands** -- run, pull, list, show, rm, cp, create, serve, ps, stop, status, models
- **Interactive REPL** -- multi-turn conversations with `bonsai run`
- **Smart model resolution** -- uses BONSAI_MODEL env var, falls back to any installed bonsai model, then any available model
- **Streaming output** -- tokens appear as they generate
- **Thinking support** -- shows a spinner during model reasoning, extracts content from thinking tokens
- **Stdin piping** -- pipe files directly into tools (`cat file.py | bonsai tool summarize`)
- **Minimal footprint** -- ~1,700 LOC, two dependencies, single binary
- **100% local** -- no telemetry, no API keys, no cloud

## Quick Start

### Prerequisites

[Ollama](https://ollama.com) must be installed and running:

```bash
# Install Ollama (macOS/Linux)
curl -fsSL https://ollama.com/install.sh | sh
ollama serve
```

### Install Bonsai

**From source:**

```bash
go install github.com/nareshnavinash/bonsai@latest
```

**From release binary:**

```bash
# macOS (Apple Silicon)
curl -L https://github.com/nareshnavinash/bonsai/releases/latest/download/bonsai-darwin-arm64 -o bonsai
chmod +x bonsai
sudo mv bonsai /usr/local/bin/
```

### First Use

```bash
# Pull a tiny but capable model (~572 MB)
bonsai pull bonsai-4b

# Ask a question
bonsai tool ask "what does defer do in Go?"

# Draft a commit message from staged changes
git add -A
bonsai tool commit

# Start an interactive chat
bonsai run
```

## Commands

### Model Management

| Command | Description |
|---------|-------------|
| `bonsai run [model] [prompt]` | Start a chat session or run a one-shot prompt |
| `bonsai pull <model>` | Download a model (supports bonsai shortnames) |
| `bonsai list` | List locally available models |
| `bonsai show <model>` | Show model details (arch, params, quantization) |
| `bonsai ps` | List currently running/loaded models |
| `bonsai stop <model>` | Unload a model from memory |
| `bonsai rm <model>` | Remove a model |
| `bonsai cp <src> <dest>` | Copy a model |
| `bonsai create <name> --from <model>` | Create a model with optional system prompt |
| `bonsai serve` | Start the Ollama server |
| `bonsai status` | Show server status and active model |
| `bonsai models` | List available Bonsai models from HuggingFace |

## Developer Tools

All tools are under `bonsai tool <name>`. They accept arguments directly or read from stdin.

### `ask` -- Ask a coding question

```bash
bonsai tool ask "how do I reverse a linked list in Python?"
```

### `commit` -- Draft a commit message

Reads staged changes via `git diff --staged` automatically:

```bash
git add -A
bonsai tool commit
# Output: feat(auth): add JWT token validation middleware
```

### `pr` -- Draft a PR description

Reads commit log from your branch's merge base with main:

```bash
bonsai tool pr
# Output: structured PR with title, summary, and test plan
```

### `diff` -- Explain staged changes

```bash
bonsai tool diff
# Output: plain-English summary of what changed and why
```

### `summarize` -- Summarize code

```bash
cat internal/chat/stream.go | bonsai tool summarize
```

### `docstring` -- Generate a docstring

```bash
cat my_function.py | bonsai tool docstring
```

### `error` -- Explain an error

```bash
bonsai tool error "panic: runtime error: index out of range [5] with length 3"

# Or pipe a stack trace:
some-command 2>&1 | bonsai tool error
```

### `regex` -- Generate a regex from description

```bash
bonsai tool regex "email address with optional subdomain"
# Output: pattern, explanation, and example matches
```

### `testdata` -- Generate test data

```bash
bonsai tool testdata "5 users with name, email, and role"
# Output: valid JSON with realistic fake data
```

### `format` -- Reformat text

```bash
cat data.csv | bonsai tool format json
cat notes.txt | bonsai tool format markdown
```

## Bonsai Models

Bonsai ships with built-in support for the **prism-ml Bonsai model family** -- true 1-bit quantized language models purpose-built for local inference.

> These models are developed by [prism-ml](https://huggingface.co/prism-ml) and available on HuggingFace:
> **[Bonsai Collection](https://huggingface.co/collections/prism-ml/bonsai-6800591046eb822fb3b82541)**

### Available Models

| Model | Parameters | Size | Pull Command |
|-------|-----------|------|-------------|
| **bonsai-8b** | 8B | ~1.2 GB | `bonsai pull bonsai-8b` |
| **bonsai-4b** | 4B | ~572 MB | `bonsai pull bonsai-4b` |
| **bonsai-1.7b** | 1.7B | ~248 MB | `bonsai pull bonsai-1.7b` |

### Why 1-Bit Models?

The prism-ml Bonsai models use **true 1-bit quantization** (not approximations with escape hatches). This means:

- **14x smaller** than FP16 equivalents
- **4-5x lower energy consumption**
- **Fast inference**: 40 tok/s on iPhone, 131 tok/s on M4 Pro, 368 tok/s on RTX 4090
- **Intelligence density**: 1.06 intelligence/GB vs 0.10 for full precision -- 10x more capability per byte

All models are in GGUF format and work natively with Ollama.

### Using Other Models

Bonsai works with any Ollama-compatible model:

```bash
bonsai pull llama3.2
bonsai run llama3.2 "explain monads"
BONSAI_MODEL=mistral bonsai tool ask "what is a goroutine?"
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `BONSAI_MODEL` | `bonsai-8b` | Model to use for tool commands |
| `OLLAMA_HOST` | `http://localhost:11434` | Ollama server address |

### Model Resolution Order

When no model is specified, Bonsai resolves in this order:

1. `BONSAI_MODEL` environment variable (if set)
2. Any locally installed model with "bonsai" in its name
3. Any locally installed model
4. Exits with a helpful message if nothing is available

### Interactive REPL Commands

In `bonsai run` interactive mode:

| Command | Description |
|---------|-------------|
| `/bye` or `/exit` | Exit the chat |
| `/set temperature <value>` | Adjust creativity (0.0-2.0) |
| `/set top_p <value>` | Adjust nucleus sampling |
| `/set num_ctx <value>` | Adjust context window size |
| `"""` | Start multi-line input (end with `"""`) |

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
