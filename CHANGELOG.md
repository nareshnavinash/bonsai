# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-04-03

### Changed
- **Removed Ollama dependency** -- bonsai now talks directly to llama.cpp server via OpenAI-compatible API (78x faster response times)
- Model downloads go directly to HuggingFace instead of through Ollama
- Models stored as GGUF files in `~/.bonsai/models/` instead of Ollama's blob store
- `serve` command now starts llama-server instead of Ollama
- `stop` command now stops the llama-server process (was: unload a specific model)
- `list` command shows local GGUF files instead of querying Ollama
- `show` command shows registry + local file info instead of Ollama model details
- `ps` command shows server process info instead of Ollama running models
- `status` command checks llama-server health instead of Ollama heartbeat
- `rm` command deletes GGUF files instead of calling Ollama delete API
- `cp` command copies GGUF files instead of calling Ollama copy API
- Environment variable `OLLAMA_HOST` replaced with `BONSAI_HOST`

### Added
- Auto-start llama-server on `bonsai run` (no need to start manually)
- Server process lifecycle management with PID tracking (`~/.bonsai/server.pid`)
- `BONSAI_HOST`, `BONSAI_PORT`, `BONSAI_THREADS`, `BONSAI_MODELS_DIR`, `LLAMA_SERVER_BIN` env vars
- Backward compatibility with legacy model paths (`~/models/bonsai-*/`)

### Removed
- `create` command (was Ollama-specific: create model with custom system prompt)
- Ollama Go SDK dependency (`github.com/ollama/ollama`) and 8 transitive dependencies
- Thinking/spinner support (no longer needed -- llama-server doesn't force thinking mode)

## [1.0.0] - 2026-04-03

### Added
- 12 model management commands: run, pull, list, show, rm, cp, create, serve, ps, stop, status, models
- Interactive REPL with multi-turn conversations
- Smart model resolution (BONSAI_MODEL -> bonsai model -> any model)
- Built-in registry for prism-ml Bonsai 1-bit quantized models (8B, 4B, 1.7B)
- Streaming output with thinking/spinner support
- Progress bar for model downloads

[2.0.0]: https://github.com/nareshnavinash/bonsai/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/nareshnavinash/bonsai/releases/tag/v1.0.0
