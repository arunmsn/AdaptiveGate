// Package policy handles rate limit decisions and quota checks.
// Dimensions: use-case-id, model-id, user-id, tenant-id.
// Sliding window, token-based and request-based. 429 with Retry-After on limit.
// Implementation: phase 2.
package policy
