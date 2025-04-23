package algo

import (
	"Go/config"
)

func AlgoLoadBalancer(algo string) string {
	selected := ""
	switch algo {
	case config.RoundRobinAlgo:
		selected = RoundRobin()

	case config.LeastConnectionAlgo:
		selected = LeastConnection()

	case config.WeightedLeastConnection:
		selected = WeightedLeastConnection()

	case config.WeightedRoundRobin:
		selected = WeightedRoundRobin()

	case config.RandomAlgo:
		selected = Random()

	default:
		selected = ""
	}
	return selected
}
