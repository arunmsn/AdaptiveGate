package main

import (
	"fmt"
	"math"
	"sort"
)

const (
	// Quality-first routing weights
	WCapability = 1.00
	WCost       = 0.18
	WLatency    = 0.12
	WFailure    = 0.10
)

type RequestContext struct {
	PromptChars       int
	ReasoningScore    float64
	CodingScore       float64
	MathScore         float64
	MultilingualScore float64

	LatencySensitive bool
	MaxCostUSDPer1M  float64

	Tenant string
}

type ModelCard struct {
	ID string

	InputUSDPer1M  float64
	OutputUSDPer1M float64

	LatencySec  float64
	FailureRate float64

	Reasoning    float64
	Coding       float64
	Math         float64
	Multilingual float64
}

type Candidate struct {
	ModelID string
	Utility float64

	Capability      float64
	CostPenalty     float64
	LatencyPenalty  float64
	FailurePenalty  float64
	NormalizedCost  float64
	NormalizedLat   float64
}

var catalog = []ModelCard{
	{
		ID: "claude-opus-4.7",

		InputUSDPer1M:  5,
		OutputUSDPer1M: 25,

		LatencySec:  1.8,
		FailureRate: 0.02,

		Reasoning:    0.98,
		Coding:       0.90,
		Math:         0.99,
		Multilingual: 0.88,
	},
	{
		ID: "gpt-5.2",

		InputUSDPer1M:  1.5,
		OutputUSDPer1M: 14,

		LatencySec:  0.6,
		FailureRate: 0.025,

		Reasoning:    0.94,
		Coding:       0.93,
		Math:         0.95,
		Multilingual: 0.86,
	},
	{
		ID: "gpt-5.3-codex",

		InputUSDPer1M:  1.75,
		OutputUSDPer1M: 14,

		LatencySec:  0.003,
		FailureRate: 0.03,

		Reasoning:    0.84,
		Coding:       0.98,
		Math:         0.88,
		Multilingual: 0.78,
	},
	{
		ID: "gemini-3-pro",

		InputUSDPer1M:  2,
		OutputUSDPer1M: 12,

		LatencySec:  30.3,
		FailureRate: 0.022,

		Reasoning:    0.96,
		Coding:       0.88,
		Math:         1.00,
		Multilingual: 0.94,
	},
	{
		ID: "deepseek-v3-0324",

		InputUSDPer1M:  0.27,
		OutputUSDPer1M: 1.10,

		LatencySec:  4,
		FailureRate: 0.035,

		Reasoning:    0.84,
		Coding:       0.78,
		Math:         0.88,
		Multilingual: 0.76,
	},
	{
		ID: "llama-4-scout",

		InputUSDPer1M:  0.11,
		OutputUSDPer1M: 0.34,

		LatencySec:  0.33,
		FailureRate: 0.04,

		Reasoning:    0.76,
		Coding:       0.70,
		Math:         0.78,
		Multilingual: 0.74,
	},
	{
		ID: "gemma-3-27b",

		InputUSDPer1M:  0.07,
		OutputUSDPer1M: 0.07,

		LatencySec:  0.72,
		FailureRate: 0.045,

		Reasoning:    0.68,
		Coding:       0.62,
		Math:         0.70,
		Multilingual: 0.72,
	},
}

func clamp01(x float64) float64 {
	switch {
	case x < 0:
		return 0
	case x > 1:
		return 1
	default:
		return x
	}
}

func minMax(values []float64) (float64, float64) {
	min := math.Inf(1)
	max := math.Inf(-1)

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return min, max
}

func normalizeMinMax(v, min, max float64) float64 {
	if math.Abs(max-min) < 1e-9 {
		return 0
	}

	return clamp01((v - min) / (max - min))
}

func estimateInputShare(promptChars int) float64 {
	if promptChars <= 0 {
		return 0.45
	}

	x := float64(promptChars) / 8000.0
	return clamp01(0.25 + 0.5*x/(x+1))
}

func blendedCost(m ModelCard, inputShare float64) float64 {
	return inputShare*m.InputUSDPer1M +
		(1-inputShare)*m.OutputUSDPer1M
}

func capabilityMatch(m ModelCard, req RequestContext) float64 {

	weights := []float64{
		req.ReasoningScore,
		req.CodingScore,
		req.MathScore,
		req.MultilingualScore,
	}

	sum := 0.0
	for _, w := range weights {
		sum += w
	}

	// neutral prior
	if sum < 1e-9 {
		return 0.75
	}

	score :=
		req.ReasoningScore*m.Reasoning +
			req.CodingScore*m.Coding +
			req.MathScore*m.Math +
			req.MultilingualScore*m.Multilingual

	return clamp01(score / sum)
}

func Route(req RequestContext) []Candidate {

	costs := make([]float64, len(catalog))
	latencies := make([]float64, len(catalog))

	inputShare := estimateInputShare(req.PromptChars)

	for i, m := range catalog {
		costs[i] = blendedCost(m, inputShare)
		latencies[i] = m.LatencySec
	}

	minCost, maxCost := minMax(costs)
	minLat, maxLat := minMax(latencies)

	latencyWeight := WLatency

	if req.LatencySensitive {
		latencyWeight *= 1.5
	}

	var candidates []Candidate

	for i, m := range catalog {

		cost := costs[i]

		if req.MaxCostUSDPer1M > 0 &&
			cost > req.MaxCostUSDPer1M {
			continue
		}

		normCost := normalizeMinMax(cost, minCost, maxCost)
		normLat := normalizeMinMax(m.LatencySec, minLat, maxLat)

		capability := capabilityMatch(m, req)

		costPenalty := WCost * normCost
		latencyPenalty := latencyWeight * normLat
		failurePenalty := WFailure * m.FailureRate

		utility :=
			WCapability*capability -
				costPenalty -
				latencyPenalty -
				failurePenalty

		candidates = append(candidates, Candidate{
			ModelID:         m.ID,
			Utility:         utility,
			Capability:      capability,
			CostPenalty:     costPenalty,
			LatencyPenalty:  latencyPenalty,
			FailurePenalty:  failurePenalty,
			NormalizedCost:  normCost,
			NormalizedLat:   normLat,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Utility >
			candidates[j].Utility
	})

	return candidates
}

func printScenario(name string, req RequestContext) {

	fmt.Println("==================================================")
	fmt.Println("SCENARIO:", name)
	fmt.Println("==================================================")

	results := Route(req)

	for i, c := range results {

		fmt.Printf(
			"%d. %s\n",
			i+1,
			c.ModelID,
		)

		fmt.Printf("   utility: %.4f\n", c.Utility)
		fmt.Printf("   capability: %.4f\n", c.Capability)
		fmt.Printf("   cost_penalty: %.4f\n", c.CostPenalty)
		fmt.Printf("   latency_penalty: %.4f\n", c.LatencyPenalty)
		fmt.Printf("   failure_penalty: %.4f\n", c.FailurePenalty)

		fmt.Println()
	}
}

func main() {

	printScenario(
		"Frontier Reasoning",
		RequestContext{
			PromptChars:      8000,
			ReasoningScore:   1.0,
			MathScore:        0.9,
			CodingScore:      0.1,
			MaxCostUSDPer1M:  100,
		},
	)

	printScenario(
		"Coding Assistant",
		RequestContext{
			PromptChars:       2000,
			CodingScore:       1.0,
			ReasoningScore:    0.4,
			LatencySensitive:  true,
			MaxCostUSDPer1M:   50,
		},
	)

	printScenario(
		"Cheap Batch Processing",
		RequestContext{
			PromptChars:      4000,
			ReasoningScore:   0.3,
			MathScore:        0.2,
			MaxCostUSDPer1M:  0.20,
		},
	)

	printScenario(
		"Multilingual Reasoning",
		RequestContext{
			PromptChars:        6000,
			ReasoningScore:     0.8,
			MultilingualScore:  1.0,
			MaxCostUSDPer1M:    100,
		},
	)
}