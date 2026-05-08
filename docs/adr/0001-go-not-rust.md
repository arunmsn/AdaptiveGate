# ADR 0001 — Go, not Rust

**Status:** accepted

## context

ixr needs to be embeddable as a library in existing services. The target audience is Go teams who want zero friction insertion into their existing stack.

## decision

Go.

## rationale

- **embed as a library** — the core value prop is `import "github.com/ixr/ixr/pkg/ixr"; ixr.Start()`. This only works if the host language matches.
- **Go is where the services are** — the primary target is teams running Go microservices. LiteLLM's Python-first design is one of the explicit problems ixr is solving.
- **no runtime overhead** — Go compiles to a single binary with no VM, no GC pauses at the latency level that matters for LLM proxying (~µs overhead).
- **stdlib is sufficient** — `net/http` handles everything we need in phase 1 without a framework.

## rejected alternatives

- **Rust** — correct choice for a standalone binary but incompatible with the embed-as-library DX. CGo bindings exist for llama.cpp in phase 2 but Rust itself isn't the right host language.
- **Python** — eliminates GC/latency problem (we'd have one) and rules out the embed path entirely.
