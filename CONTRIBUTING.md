# Contributing to ixr

## DCO

All contributions must be signed off under the [Developer Certificate of Origin](https://developercertificate.org/).

Add a `Signed-off-by` trailer to every commit:

```
git commit -s -m "your message"
```

## PR guidelines

- **one feature per PR** — small diffs, fast review, clear changelog entries
- **one binary to run** — no new required infrastructure without an ADR
- **dependency rule is sacred** — `internal/domain` must never import `internal/adapters`; CI enforces this
- **table-driven tests** — every new package needs a `_test.go`
- **comments explain why, not what** — if the what isn't obvious, rename the function

## Changes to `pkg/`

`pkg/` is semver-governed public API. Breaking changes require an RFC filed as a GitHub issue before any code is written.

## ADRs

Significant architectural decisions go in `docs/adr/` as numbered markdown files. See existing ADRs for the format.

## CI

All PRs must pass:

```
go test ./...
go vet ./...
staticcheck ./...
make check-deps
```
