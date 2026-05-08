// Package scoring contains the adaptive scoring engine.
// v1: deterministic weighted scoring driven by per-intent weights from the policy store.
// v2: bandit algorithms (epsilon-greedy and UCB) that learn from real traffic.
package scoring
