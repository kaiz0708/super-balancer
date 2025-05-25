package algo

import (
	"Go/algo/custom"
	algo_default "Go/algo/defaults"
	"Go/config"
	"net/http"
)

var AlgoCurrent string

func ChooseAlgorithm(state string, r *http.Request) string {
	selected := ""
	if !config.ConfigSystem.SmartMode {
		selected = AlgoLoadBalancer(config.ConfigSystem.Algorithm, r)
		return selected
	}
	switch state {
	case config.ManyFailed:
		config.ConfigSystem.Algorithm = config.WeightedSuccessRateBalancer
		selected = custom.WeightedSuccessRateBalancer()
	case config.HighLatency:
		config.ConfigSystem.Algorithm = config.LowLatencyWeightedBalancer
		selected = custom.LowLatencyWeightedBalancer()
	default:
		selected = AlgoLoadBalancer(config.ConfigSystem.Algorithm, r)
	}
	return selected
}

func AlgoLoadBalancer(algo string, r *http.Request) string {
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

	case config.WeightedRandom:
		selected = algo_default.WeightedRandom()

	case config.IpHash:
		selected = algo_default.IpHash(r)

	default:
		selected = algo_default.RoundRobin()
	}
	return selected
}
