package intent

// Intent constants define the supported intent taxonomy (v1).
// Intents are extensible via config without changing code.
const (
	IntentReasoning      = "reasoning"      // multi-step logic, math, code — biases accuracy over cost
	IntentSummarization  = "summarization"  // condensing long input — biases cost efficiency, speed
	IntentExtraction     = "extraction"     // structured output from text — biases reliability, consistency
	IntentGeneration     = "generation"     // creative / open-ended output — biases quality, model capability
	IntentClassification = "classification" // label or categorize input — biases speed, cost
	IntentEmbedding      = "embedding"      // vector representation — biases specialized embedding models
)

// DefaultWeights holds the v1 per-intent scoring weights (w1=cost, w2=latency, w3=reliability).
// These are overridden by the policy store at route time.
var DefaultWeights = map[string][3]float64{
	IntentReasoning:      {0.2, 0.3, 0.5},
	IntentSummarization:  {0.5, 0.3, 0.2},
	IntentExtraction:     {0.3, 0.2, 0.5},
	IntentGeneration:     {0.2, 0.2, 0.6},
	IntentClassification: {0.5, 0.4, 0.1},
}
