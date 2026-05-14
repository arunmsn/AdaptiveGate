package main

import (
	"math"
	"testing"
)

func TestClamp01(t *testing.T) {

	if clamp01(-1) != 0 {
		t.Fatal("negative clamp failed")
	}

	if clamp01(2) != 1 {
		t.Fatal("upper clamp failed")
	}

	if clamp01(0.5) != 0.5 {
		t.Fatal("identity clamp failed")
	}
}

func TestNormalizeMinMax(t *testing.T) {

	got := normalizeMinMax(5, 0, 10)

	if math.Abs(got-0.5) > 1e-9 {
		t.Fatalf("expected 0.5 got %f", got)
	}
}

func TestCapabilityMatch(t *testing.T) {

	req := RequestContext{
		CodingScore: 1.0,
	}

	codingModel := ModelCard{
		Coding: 0.95,
	}

	weakModel := ModelCard{
		Coding: 0.20,
	}

	if capabilityMatch(codingModel, req) <=
		capabilityMatch(weakModel, req) {

		t.Fatal("coding model should score higher")
	}
}

func TestReasoningRouting(t *testing.T) {

	req := RequestContext{
		ReasoningScore:  1.0,
		MathScore:       0.9,
		MaxCostUSDPer1M: 100,
	}

	results := Route(req)

	if len(results) == 0 {
		t.Fatal("no routing results")
	}

	top := results[0]

	valid :=
		top.ModelID == "claude-opus-4.7" ||
			top.ModelID == "gpt-5.2"

	if !valid {
		t.Fatalf(
			"unexpected reasoning winner: %s",
			top.ModelID,
		)
	}
}

func TestCodingRouting(t *testing.T) {

	req := RequestContext{
		CodingScore:      1.0,
		LatencySensitive: true,
		MaxCostUSDPer1M:  50,
	}

	results := Route(req)

	if results[0].ModelID != "gpt-5.3-codex" {
		t.Fatalf(
			"expected gpt-5.3-codex got %s",
			results[0].ModelID,
		)
	}
}

func TestBudgetFiltering(t *testing.T) {

	req := RequestContext{
		ReasoningScore:  0.5,
		MaxCostUSDPer1M: 0.10,
	}

	results := Route(req)

	if len(results) == 0 {
		t.Fatal("budget filter removed everything")
	}

	if results[0].ModelID != "gemma-3-27b" {
		t.Fatalf(
			"expected cheapest model got %s",
			results[0].ModelID,
		)
	}
}

func TestZeroVectorRequest(t *testing.T) {

	req := RequestContext{}

	results := Route(req)

	if len(results) == 0 {
		t.Fatal("expected fallback routing")
	}
}

func TestLatencySensitiveRouting(t *testing.T) {

	req := RequestContext{
		CodingScore:      0.8,
		LatencySensitive: true,
		MaxCostUSDPer1M:  100,
	}

	results := Route(req)

	top := results[0]

	if top.NormalizedLat > 0.5 {
		t.Fatal("latency-sensitive routing picked slow model")
	}
}

func TestSortedUtilities(t *testing.T) {

	req := RequestContext{
		ReasoningScore: 1.0,
	}

	results := Route(req)

	for i := 1; i < len(results); i++ {

		if results[i].Utility >
			results[i-1].Utility {

			t.Fatal("results not sorted descending")
		}
	}
}
