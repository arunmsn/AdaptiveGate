package main

import (
	"fmt"
	"math"
	"sort"
)

const (
	// Infrastructure-first routing weights
	WCost     = 0.45
	WLatency  = 0.35
	WFailure  = 0.20
)

type ModelCard struct {
	ID string

	// USD per 1M tokens
	InputUSDPer1M  float64
	OutputUSDPer1M float64

	// Representative latency
	LatencySec float64

	// Failure probability
	FailureRate float64
}

type Candidate struct {
	ModelID string

	Penalty float64

	NormalizedCost    float64
	NormalizedLatency float64
	FailureRate       float64

	CostContribution    float64
	LatencyContribution float64
	FailureContribution float64
}

var catalog = []ModelCard{
	{
		ID: "claude-opus-4.7",

		InputUSDPer1M:  5,
		OutputUSDPer1M: 25,

		LatencySec:  1.8,
		FailureRate: 0.02,
	},
	{
		ID: "gpt-5.2",

		InputUSDPer1M:  1.5,
		OutputUSDPer1M: 14,

		LatencySec:  0.6,
		FailureRate: 0.025,
	},
	{
		ID: "gpt-5.3-codex",

		InputUSDPer1M:  1.75,
		OutputUSDPer1M: 14,

		LatencySec:  0.003,
		FailureRate: 0.03,
	},
	{
		ID: "gemini-3-pro",

		InputUSDPer1M:  2,
		OutputUSDPer1M: 12,

		LatencySec:  30.3,
		FailureRate: 0.022,
	},
	{
		ID: "deepseek-v3-0324",

		InputUSDPer1M:  0.27,
		OutputUSDPer1M: 1.10,

		LatencySec:  4,
		FailureRate: 0.035,
	},
	{
		ID: "llama-4-scout",

		InputUSDPer1M:  0.11,
		OutputUSDPer1M: 0.34,

		LatencySec:  0.33,
		FailureRate: 0.04,
	},
	{
		ID: "gemma-3-27b",

		InputUSDPer1M:  0.07,
		OutputUSDPer1M: 0.07,

		LatencySec:  0.72,
		FailureRate: 0.045,
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

// Simple blended cost estimate
func blendedCost(m ModelCard) float64 {

	inputWeight := 0.5
	outputWeight := 0.5

	return inputWeight*m.InputUSDPer1M +
		outputWeight*m.OutputUSDPer1M
}

func Route() []Candidate {

	costs := make([]float64, len(catalog))
	latencies := make([]float64, len(catalog))

	for i, m := range catalog {

		costs[i] = blendedCost(m)
		latencies[i] = m.LatencySec
	}

	minCost, maxCost := minMax(costs)
	minLat, maxLat := minMax(latencies)

	var candidates []Candidate

	for i, m := range catalog {

		normCost :=
			normalizeMinMax(costs[i], minCost, maxCost)

		normLatency :=
			normalizeMinMax(m.LatencySec, minLat, maxLat)

		costContribution :=
			WCost * normCost

		latencyContribution :=
			WLatency * normLatency

		failureContribution :=
			WFailure * m.FailureRate

		penalty :=
			costContribution +
				latencyContribution +
				failureContribution

		candidates = append(candidates, Candidate{
			ModelID: m.ID,

			Penalty: penalty,

			NormalizedCost:    normCost,
			NormalizedLatency: normLatency,
			FailureRate:       m.FailureRate,

			CostContribution:    costContribution,
			LatencyContribution: latencyContribution,
			FailureContribution: failureContribution,
		})
	}

	// lower penalty is better
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Penalty <
			candidates[j].Penalty
	})

	return candidates
}

func printResults() {

	fmt.Println("==================================================")
	fmt.Println("INFRASTRUCTURE-FIRST ROUTING")
	fmt.Println("==================================================")

	results := Route()

	for i, c := range results {

		fmt.Printf("%d. %s\n", i+1, c.ModelID)

		fmt.Printf("   total_penalty: %.4f\n", c.Penalty)

		fmt.Printf("   normalized_cost: %.4f\n",
			c.NormalizedCost)

		fmt.Printf("   normalized_latency: %.4f\n",
			c.NormalizedLatency)

		fmt.Printf("   failure_rate: %.4f\n",
			c.FailureRate)

		fmt.Printf("   cost_contribution: %.4f\n",
			c.CostContribution)

		fmt.Printf("   latency_contribution: %.4f\n",
			c.LatencyContribution)

		fmt.Printf("   failure_contribution: %.4f\n",
			c.FailureContribution)

		fmt.Println()
	}
}

func main() {
	printResults()
}

