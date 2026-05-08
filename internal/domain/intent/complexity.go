package intent

// ComplexityBucket classifies a request's complexity as a routing signal.
// v1 uses a heuristic: prompt length + presence of code blocks + chain-of-thought markers.
// Implementation: phase 2.
const (
	ComplexityLow    = "low"
	ComplexityMedium = "medium"
	ComplexityHigh   = "high"
)
