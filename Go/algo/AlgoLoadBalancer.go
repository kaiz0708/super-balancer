package algo

import (
	"Go/algo/custom"
	algo_default "Go/algo/defaults"
	"Go/config"
)

var AlgoCurrent string

func ChooseAlgorithm(state string) string {
	selected := ""
	switch state {
	case "ManyFailed":
		AlgoCurrent = "ManyFailed"
		selected = custom.CustomManyFailed()
	case "HighLatency":
		AlgoCurrent = "HighLatency"
		selected = custom.CustomHighLatency()
	default:
		AlgoCurrent = config.LoadBalancerDefault
		selected = AlgoLoadBalancer(config.LoadBalancerDefault)
	}
	return selected
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
