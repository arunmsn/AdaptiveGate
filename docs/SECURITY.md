# ixr security

## threat model

ixr sits between your service and LLM providers. The primary surfaces:

1. **Ingress** — unauthenticated in phase 1 (bind to localhost); auth is phase 2
2. **Provider credentials** — API keys in env vars; secret manager support in phase 2
3. **Plugin bus** — plugins run in-process; a malicious plugin has full access
4. **Supply chain** — minimal deps, signed releases, SBOM on every release

## reporting vulnerabilities

See [SECURITY.md](../SECURITY.md) at the repo root.

## phase 2 additions

- Ingress auth: API key, mTLS, JWT verification
- Per-key scoping (this key can only call gpt-3.5, this key has a $100/day cap)
- Secrets rotation without restart (Vault, AWS Secrets Manager, GCP Secret Manager)
- PII-block plugin (reference implementation)
