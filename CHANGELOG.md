# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Repo scaffold: module, folder structure, public interfaces, CI skeleton
- `pkg/schema` — `CallEvent`, `RequestEnvelope`, `ResponseEnvelope`, `TelemetryRecord` types
- `pkg/plugin` — `EventConsumer` interface
- `pkg/provider` — `Provider` interface
- `pkg/bus` — `Bus` interface
- `pkg/ixr` — `Start()` facade (stub)
- Apache 2.0 license
- GitHub Actions: test, release, govulncheck workflows
