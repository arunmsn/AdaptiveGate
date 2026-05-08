// Package modelperf stores and retrieves per-(model, intent) performance statistics.
// postgres is the source of truth; redis is the hot cache the scoring engine reads.
package modelperf
