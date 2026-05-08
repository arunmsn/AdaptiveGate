# ADR 0003 — Event bus shape

**Status:** accepted

## context

ixr needs to deliver `CallEvent`s to plugins without blocking the request path. The bus implementation will evolve from in-memory to NATS/Kafka in phase 2.

## decision

Define a `pkg/bus.Bus` interface with two methods: `Publish` and `Subscribe`. Phase 1 implements it as a buffered in-process Go channel.

## rationale

- **interface at the seam** — callers (`internal/app`) depend on `pkg/bus.Bus`, not on the concrete implementation. Swapping in-memory for NATS changes one line in `cmd/ixr/main.go`.
- **non-blocking Publish** — the channel is buffered. A slow plugin never adds latency to the caller's request.
- **no external dependencies in phase 1** — no NATS, no Kafka, no Kafka client libs. The in-memory implementation has zero deps.

## rejected alternatives

- **Direct function calls to plugins** — couples the app layer to every plugin. Adding a plugin means touching core code.
- **gRPC streaming** — overkill for phase 1, and adds a network hop for in-process consumers.
- **Redis pub/sub** — requires Redis in phase 1 which violates the zero-infra requirement.
