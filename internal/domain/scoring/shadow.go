package scoring

// shadow orchestrates shadow routing: the primary model serves the caller while
// a shadow model processes the same request in the background for offline comparison.
// Shadow routing is opt-in per use-case, configured via the policy store.
// Implementation: phase 2.
