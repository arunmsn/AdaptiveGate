package scoring

// reward computes: α*(1/latency_ms) + β*(1-cost_per_token) + γ*success_rate + δ*quality_score.
// α, β, γ, δ are learned per intent by the bandit algorithm.
// Implementation: phase 2.
