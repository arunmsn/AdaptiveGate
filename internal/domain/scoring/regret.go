package scoring

// regret tracks cumulative regret = sum(optimal_reward - chosen_reward) across all requests.
// Lower cumulative regret means the algorithm is learning faster.
// This is the north star metric for v2 routing quality.
// Implementation: phase 2.
