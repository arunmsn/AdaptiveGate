# ADR 0004 — Intent as headers, not a new API schema

**Status:** accepted

## context

Callers who want smarter routing need a way to express intent and constraints to ixr. The question is where that signal lives.

## decision

Intent and constraints are expressed as HTTP headers (`X-IXR-Intent`, `X-IXR-Max-Cost`, etc.) or as an optional `x_ixr` field in the existing JSON body.

## rationale

- **zero breaking change** — existing callers that don't set the headers continue to work exactly as before. ixr routes via model-prefix rules.
- **no new API schema** — callers don't need to change their SDK usage. They add headers. This is the lowest friction possible.
- **`model: "auto"` as the trigger** — callers who want full auto-routing set `"model": "auto"`. This is a recognized OpenAI pattern (they use it for model selection too), so it's intuitive.
- **additive, not required** — the intent layer is opt-in. This keeps the phase 1 DX simple.

## rejected alternatives

- **New endpoint** (`POST /v1/chat/completions/smart`) — breaks the zero-refactor guarantee. Callers would need to change their SDK base URL pattern.
- **Required body field** — breaks existing callers that don't populate it.
