// Package routing picks a provider + model for each request given a ParsedRequest.
// It orchestrates: filter → score → select → build fallback chain.
// Pure domain logic — no HTTP, no Redis, no external dependencies.
package routing

// RoutingDecision is the output of the routing engine for a single request.
type RoutingDecision struct {
	Provider      string
	Model         string
	FallbackChain []Candidate
}

// Candidate is a scored model eligible for routing.
type Candidate struct {
	Provider string
	Model    string
	Score    float64
}
