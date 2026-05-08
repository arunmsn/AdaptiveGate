# Writing ixr provider adapters

A provider adapter implements `pkg/provider.Provider` and translates ixr's canonical schema to and from a specific LLM's wire format.

## the interface

```go
type Provider interface {
    Name() string
    Chat(ctx context.Context, req *schema.RequestEnvelope) (*schema.ResponseEnvelope, error)
}
```

## adding a new provider

1. Create `internal/adapters/providers/<name>/adapter.go`
2. Implement `Provider` — `Name()` returns the config key (e.g. `"bedrock"`)
3. Create `translate.go` — the schema ↔ wire format conversions
4. Create `adapter_test.go` with recorded fixtures (no live API key needed in CI)
5. Register in the provider registry (phase 1 day 5)

## model prefix routing

Phase 1 routes by model name prefix:
- `gpt-*` → openai
- `claude-*` → anthropic

Phase 2 routing is driven by the scoring engine.

## recorded fixtures

Use a `httptest.Server` or a response-recording transport to capture real API responses once. Tests replay from fixtures — no live keys needed in CI.

See `internal/adapters/providers/openai/adapter_test.go` for the pattern.
