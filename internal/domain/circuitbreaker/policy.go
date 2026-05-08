package circuitbreaker

// policy holds circuit breaker thresholds and probe logic.
// Defaults: open when success_rate < 0.90 over 2 min; half-open after 30s.
// All thresholds are configurable per provider.
// Implementation: phase 2.
