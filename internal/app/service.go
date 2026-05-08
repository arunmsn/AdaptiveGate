// Package app is the application / orchestration layer.
// It coordinates domain pieces in the right order but holds no business logic itself.
// Pipeline: request → intent parser → scoring engine → plugins → provider → plugins → response.
package app
