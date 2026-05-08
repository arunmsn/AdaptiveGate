# ixr architecture

ixr follows a clean layered architecture: **api boundaries → application orchestration → domain logic → infrastructure adapters**. Each layer has one responsibility and depends only inward.

## layers

```
cmd/              presentation — wires dependencies and starts the app
pkg/              public api — semver-governed, third parties import these
internal/
  ingress/        transport — HTTP, parses OpenAI-shaped requests
  app/            orchestration — coordinates domain pieces, no business logic
  domain/         pure logic — no HTTP, no providers, no infra; fully unit-testable
  adapters/       infrastructure — everything that talks to a network or disk
plugins/          reference plugins — each is a separate go module
```

## dependency rule

```
cmd → pkg, internal/app
           ↓
      internal/app → internal/domain, internal/adapters, pkg/*
                         ↓                    ↓
                 internal/domain      internal/adapters → pkg/*
                   (nothing internal)   (pkg interfaces only)
```

**Code only depends inward (toward `domain/`), never outward.** CI enforces this with `go list` checks via `make check-deps`.

Violations:
- If `internal/domain/routing` imports `internal/adapters/providers/openai` → architecture is broken
- If `internal/adapters/providers/openai` imports `internal/adapters/providers/anthropic` → architecture is broken

## request pipeline

```
POST /v1/chat/completions
        ↓
   internal/ingress        parse, normalize → RequestEnvelope
        ↓
   internal/app            orchestrate stages
        ↓
   domain/intent           parse X-IXR-* headers → ParsedRequest
        ↓
   domain/scoring          pick model, build fallback chain → RoutingDecision
        ↓
   plugin manager          pre-call plugins (phase 2)
        ↓
   adapters/providers      call OpenAI / Anthropic → ResponseEnvelope
        ↓
   plugin manager          publish CallEvent to bus (async, non-blocking)
        ↓
   internal/ingress        shape response as OpenAI format → caller
```

## the bus

Every `CallEvent` is published to the bus after the provider call completes. Bus delivery is:
- **non-blocking** — a slow plugin never affects request latency
- **in-process** for phase 1 (a buffered Go channel)
- **swappable** in phase 2 (NATS, Kafka, Kinesis — same `pkg/bus.Bus` interface)

## design principles

See the full list in the PRD. Key ones:

1. **stdlib first** — no gin, echo, or fiber
2. **interfaces at the seams** — providers, plugins, bus, stores are interfaces; data is plain structs
3. **no init()** — all wiring is explicit and ordered in `cmd/ixr/main.go` or `pkg/ixr.Start()`
4. **no global state** — every component takes its dependencies; testable by construction
5. **public api frozen at v1** — `pkg/` follows semver; `internal/` is fair game
