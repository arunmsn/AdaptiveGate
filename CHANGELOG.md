# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-08

### Added
- `pkg/ixr` — `Start(opts ...Option) error` one-line entry point; `WithPort`, `WithConfigFile` options
- `pkg/schema` — `CallEvent`, `RequestEnvelope`, `ResponseEnvelope`, `Message`, `Choice`, `Usage`, `CostBreakdown`, `ToolCall` public types
- `pkg/plugin` — `EventConsumer` interface for zero-fork extensibility
- `pkg/provider` — `Provider` interface (`Name`, `Chat`)
- `pkg/bus` — `Bus` interface (`Publish`, `Subscribe`)
- OpenAI provider — full `POST /v1/chat/completions` passthrough, OpenAI-compatible response shape
- Anthropic provider — Messages API integration, system-message lifting, stop-reason normalisation
- Model-prefix router — `gpt-*` / `o1` / `o3` → OpenAI; `claude-*` → Anthropic
- In-memory event bus — buffered channel, non-blocking publish, panic-safe plugin dispatch
- Plugin manager — registers `EventConsumer` plugins at startup
- Audit-log plugin — emits every `CallEvent` as a JSON line to stdout
- Config loader — `ixr.yaml` with `${ENV_VAR}` interpolation, auto-discovery, env-var override
- `cmd/ixr` binary — `--config` and `--port` flags
- `Dockerfile` — multi-stage scratch image, `linux/amd64` + `linux/arm64`
- Table-driven tests — translators, adapters, chat handler, config loader (no live API keys required)
- GitHub Actions — test (go vet + staticcheck + race detector), release (multi-arch binaries, cosign, syft SBOM, ghcr.io image), govulncheck
- Apache 2.0 license

[0.1.0]: https://github.com/YashVishwas/ixr/releases/tag/v0.1.0
