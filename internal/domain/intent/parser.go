// Package intent parses and enriches intent + constraints from the request.
// It maps X-IXR-* headers and x_ixr body fields into a ParsedRequest that
// the routing engine uses to pick a model.
package intent

// ParsedRequest is the enriched representation of an inbound request
// after intent and constraint extraction.
type ParsedRequest struct {
	Intent           string   // one of the taxonomy constants
	MaxCostUSD       *float64 // hard ceiling; nil = unconstrained
	MaxLatencyMS     *int     // hard ceiling; nil = unconstrained
	MinQuality       *float64 // 0.0–1.0; nil = unconstrained
	TokenEstimate    int      // derived from prompt token count
	ComplexityBucket string   // "low" | "medium" | "high" (heuristic)
}
