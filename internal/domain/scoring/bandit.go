package scoring

// bandit implements epsilon-greedy and UCB adaptive routing algorithms.
// Both run in shadow mode before either becomes the live algorithm.
// The better one is chosen by minimizing cumulative regret on real traffic.
// Implementation: phase 2.
