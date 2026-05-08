# ADR 0002 — pkg/ vs internal/

**Status:** accepted

## context

Third parties need stable types to build intelligence layers on top of ixr. At the same time, ixr's implementation needs to evolve freely.

## decision

Split the codebase into two zones:
- `pkg/` — public, semver-governed, third parties import these
- `internal/` — private, can change at any time

## rationale

Go's `internal/` package visibility rules enforce this at the compiler level. No accidental public API exposure.

`pkg/` contains exactly:
- `pkg/schema` — the data contracts (`CallEvent`, envelopes, `TelemetryRecord`)
- `pkg/plugin` — the `EventConsumer` interface
- `pkg/provider` — the `Provider` interface
- `pkg/bus` — the `Bus` interface
- `pkg/ixr` — the `Start()` facade

Everything else is `internal/`. This is a small surface to stabilize and a large surface to evolve.

## consequence

A breaking change to `pkg/schema.CallEvent` requires a semver major bump and a migration guide. This is the intended constraint — it forces deliberate API design.
