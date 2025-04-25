package algo

import (
	"Go/algo/custom"
	algo_default "Go/algo/defaults"
	"Go/config"
)

func ChooseAlgorithm(state string) {
	switch state {
	case "ManyFailed":
		config.LoadBalancerDefault = custom.CustomManyFailed()
	case "HighLatency":
		config.LoadBalancerDefault = custom.CustomHighLatency()
	default:
		config.LoadBalancerDefault = algo_default.WeightedRoundRobin()
	}
}

func AlgoLoadBalancer(algo string) string {
	selected := ""
	switch algo {
	case config.RoundRobinAlgo:
		selected = algo_default.RoundRobin()

	case config.LeastConnectionAlgo:
		selected = algo_default.LeastConnection()

	case config.WeightedLeastConnection:
		selected = algo_default.WeightedLeastConnection()

	case config.WeightedRoundRobin:
		selected = algo_default.WeightedRoundRobin()

	case config.RandomAlgo:
		selected = algo_default.Random()

	default:
		selected = algo_default.RoundRobin()
	}
	return selected
}
