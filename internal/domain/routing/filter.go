package routing

// filter removes candidate models that violate hard constraints:
// max_cost_usd, max_latency_ms (p95), and circuit-breaker exclusions.
// Implementation: phase 2 (scoring engine).
