package modelperf

// redis is the hot cache for scoring engine reads (target: < 1ms).
// All data is pre-computed and written here by the telemetry pipeline.
// Implementation: phase 2.
