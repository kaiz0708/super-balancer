package algo

import (
	"Go/config"
)

func ChooseAlgorithm(state string) {
	switch state {
	case "AllFailed":
		config.LoadBalancerDefault = Random()
	case "ManyFailed":
		config.LoadBalancerDefault = WeightedLeastConnection()
	case "HighLatency":
		config.LoadBalancerDefault = WeightedRoundRobin()
	default:
		config.LoadBalancerDefault = RoundRobin()
	}
}

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
		selected = RoundRobin()
	}
	return selected
}
