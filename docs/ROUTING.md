# ixr routing

## phase 1 — model prefix routing

Simple, deterministic:
- `gpt-*` → openai
- `claude-*` → anthropic

No config required beyond provider API keys.

## phase 2 — intent-based scoring engine

Opt-in via request headers:

```
X-IXR-Intent: reasoning
X-IXR-Max-Cost: 0.01
X-IXR-Max-Latency: 1500
X-IXR-Quality: high
```

Or via request body:

```json
{
  "model": "auto",
  "x_ixr": {
    "intent": "reasoning",
    "constraints": { "max_cost_usd": 0.01, "max_latency_ms": 1500 }
  }
}
```

When `model` is `"auto"` or intent headers are present, the scoring engine picks the model. **Zero breaking change** for callers that don't opt in.

## scoring algorithm (v1 — deterministic)

1. **Filter** — remove models violating hard constraints (cost, latency p95, circuit-open)
2. **Score** — `score(model) = w1*cost + w2*latency + w3*(1-success_rate)` (lower = better)
3. **Select** — pick lowest score
4. **Fallback chain** — next 2 lowest scores become the fallback chain

Weights `w1, w2, w3` are per-intent and configurable in the policy store.

## intent taxonomy

| intent | biases toward |
|--------|--------------|
| `reasoning` | accuracy over cost |
| `summarization` | cost efficiency, speed |
| `extraction` | reliability, consistency |
| `generation` | quality, model capability |
| `classification` | speed, cost |
| `embedding` | specialized embedding models |

## v2 — adaptive routing

See [ADAPTIVE.md](ADAPTIVE.md).
