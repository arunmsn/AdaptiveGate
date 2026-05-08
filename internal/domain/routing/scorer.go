package routing

// scorer computes score(model) = w1*normalized_cost + w2*normalized_latency + w3*(1-success_rate).
// Lower score = better candidate.
// Weights are per-intent and loaded from the policy store at route time.
// Implementation: phase 2 (scoring engine).
