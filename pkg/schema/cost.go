package schema

// CostBreakdown holds the USD cost components for a single LLM call.
type CostBreakdown struct {
	InputUSD  float64 `json:"input_usd"`
	OutputUSD float64 `json:"output_usd"`
	TotalUSD  float64 `json:"total_usd"`
}
