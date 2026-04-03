# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-04-03

### Added
- 12 model management commands: run, pull, list, show, rm, cp, create, serve, ps, stop, status, models
- 10 AI-powered developer tools: ask, summarize, docstring, error, commit, diff, pr, regex, testdata, format
- Interactive REPL with multi-turn conversations
- Smart model resolution (BONSAI_MODEL -> bonsai model -> any model)
- Git-aware tools: commit reads staged diffs, pr reads branch history
- Streaming output with thinking/spinner support
- Built-in registry for prism-ml Bonsai 1-bit quantized models (8B, 4B, 1.7B)
- Progress bar for model downloads
- Stdin piping support for summarize, docstring, error, and format

[1.0.0]: https://github.com/nareshnavinash/bonsai/releases/tag/v1.0.0
