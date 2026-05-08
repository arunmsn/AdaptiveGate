# ADR 0005 — Scoring engine in domain/, not app/

**Status:** accepted

## context

The scoring engine is the brain of ixr's routing. It needs to be placed in the layer hierarchy.

## decision

The scoring engine lives in `internal/domain/scoring/`, not in `internal/app/`.

## rationale

The scoring engine is **pure business logic** — given a `ParsedRequest` and a set of model statistics, it returns a `RoutingDecision`. It has no side effects, no HTTP calls, no database queries. It reads from interfaces (`ModelPerfStore`, `PolicyStore`) that are defined in `pkg/` and implemented in `internal/adapters/store/`.

Putting it in `domain/` means:
- It can be unit-tested without any infrastructure
- It never imports from `internal/adapters/`
- The dependency rule is maintained

If it were in `app/`, it would be tempting to give it direct Redis access, which would couple it to infrastructure and make it impossible to test without a Redis instance.

## consequence

The Redis adapter is in `internal/adapters/store/redis.go`. The domain scoring engine reads from a `ModelPerfStore` interface. The wiring happens in `cmd/ixr/main.go`. The domain stays clean.
