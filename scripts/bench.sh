#!/usr/bin/env bash
# bench.sh — run benchmarks and report added latency
# usage: ./scripts/bench.sh
set -euo pipefail

echo "Running benchmarks..."
go test -bench=. -benchmem -count=5 ./... | tee bench.txt

echo ""
echo "Benchmark results saved to bench.txt"
