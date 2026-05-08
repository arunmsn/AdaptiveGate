// Package circuitbreaker watches model health and excludes degraded models from routing.
// States: closed (healthy) → open (excluded) → half-open (probe) → closed.
// Circuit state is stored in Redis so all ixr instances see it immediately.
package circuitbreaker
