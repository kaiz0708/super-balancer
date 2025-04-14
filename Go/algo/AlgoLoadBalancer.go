package algo

import (
	"os"
)

func AlgoLoadBalancer(algo string) string {
	selected := ""
	switch algo {
	case os.Getenv("ROUND_ROBIN"):
		selected = RoundRobin()

	case os.Getenv("LEAST_CONNECTION"):
		selected = LeastConnection()

	case os.Getenv("WEIGHTED_LEAST_CONNECTION"):
		selected = WeightedLeastConnection()

	case os.Getenv("WEIGHTED_ROUND_ROBIN"):
		selected = WeightedRoundRobin()

	case os.Getenv("RANDOM"):
		selected = Random()

	default:
		selected = ""
	}
	return selected
}
