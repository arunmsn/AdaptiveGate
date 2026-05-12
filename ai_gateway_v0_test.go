package main

import (
	"math"
	"testing"
)

func TestClamp01(t *testing.T) {

	tests := []struct {
		input    float64
		expected float64
	}{
		{-1, 0},
		{0, 0},
		{0.5, 0.5},
		{1, 1},
		{2, 1},
	}

	for _, tt := range tests {

		got := clamp01(tt.input)

		if got != tt.expected {

			t.Fatalf(
				"clamp01(%f) = %f, expected %f",
				tt.input,
				got,
				tt.expected,
			)
		}
	}
}

func TestNormalizeMinMax(t *testing.T) {

	got := normalizeMinMax(5, 0, 10)

	if math.Abs(got-0.5) > 1e-9 {

		t.Fatalf(
			"expected 0.5 got %f",
			got,
		)
	}
}

func TestNormalizeMinMaxEqualBounds(t *testing.T) {

	got := normalizeMinMax(5, 5, 5)

	if got != 0 {

		t.Fatalf(
			"expected 0 got %f",
			got,
		)
	}
}

func TestBlendedCost(t *testing.T) {

	model := ModelCard{
		InputUSDPer1M:  2,
		OutputUSDPer1M: 10,
	}

	got := blendedCost(model)

	expected := 6.0

	if math.Abs(got-expected) > 1e-9 {

		t.Fatalf(
			"expected %f got %f",
			expected,
			got,
		)
	}
}

func TestRouteReturnsCandidates(t *testing.T) {

	results := Route()

	if len(results) == 0 {
		t.Fatal("expected candidates")
	}
}

func TestResultsSortedAscendingPenalty(t *testing.T) {

	results := Route()

	for i := 1; i < len(results); i++ {

		if results[i].Penalty <
			results[i-1].Penalty {

			t.Fatal(
				"results not sorted by ascending penalty",
			)
		}
	}
}

func TestCheapModelsRankHigher(t *testing.T) {

	results := Route()

	top := results[0]

	valid :=
		top.ModelID == "gemma-3-27b" ||
			top.ModelID == "llama-4-scout"

	if !valid {

		t.Fatalf(
			"unexpected top infra model: %s",
			top.ModelID,
		)
	}
}

func TestGeminiPenalizedForLatency(t *testing.T) {

	results := Route()

	var gemini Candidate

	found := false

	for _, r := range results {

		if r.ModelID == "gemini-3-pro" {
			gemini = r
			found = true
			break
		}
	}

	if !found {
		t.Fatal("gemini model not found")
	}

	if gemini.NormalizedLatency < 0.9 {

		t.Fatalf(
			"expected very high latency penalty got %f",
			gemini.NormalizedLatency,
		)
	}
}

func TestClaudePenalizedForCost(t *testing.T) {

	results := Route()

	var claude Candidate

	found := false

	for _, r := range results {

		if r.ModelID == "claude-opus-4.7" {
			claude = r
			found = true
			break
		}
	}

	if !found {
		t.Fatal("claude model not found")
	}

	if claude.NormalizedCost < 0.8 {

		t.Fatalf(
			"expected high normalized cost got %f",
			claude.NormalizedCost,
		)
	}
}

func TestFailureContributionPositive(t *testing.T) {

	results := Route()

	for _, r := range results {

		if r.FailureContribution < 0 {

			t.Fatal(
				"failure contribution should never be negative",
			)
		}
	}
}

func TestPenaltyNonNegative(t *testing.T) {

	results := Route()

	for _, r := range results {

		if r.Penalty < 0 {

			t.Fatal(
				"penalty should never be negative",
			)
		}
	}
}