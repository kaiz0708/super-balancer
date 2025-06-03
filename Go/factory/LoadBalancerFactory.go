package factory

import (
	"Go/algo/custom"
	algo_default "Go/algo/defaults"
	"net/http"
)

type LoadBalancerStrategy interface {
	SelectServer(r *http.Request) string
}

type LoadBalancerFactory struct {
	strategies map[string]LoadBalancerStrategy
}

func NewLoadBalancerFactory() *LoadBalancerFactory {
	return &LoadBalancerFactory{
		strategies: map[string]LoadBalancerStrategy{
			"ROUND_ROBIN":                    &algo_default.RoundRobinStrategy{},
			"LEAST_CONNECTION":               &algo_default.LeastConnectionStrategy{},
			"IP_HASH":                        &algo_default.IpHashStrategy{},
			"RANDOM":                         &algo_default.RandomStrategy{},
			"WEIGHTED_LEAST_CONNECTION":      &algo_default.WeightedLeastConnectionStrategy{},
			"WEIGHTED_ROUND_ROBIN":           &algo_default.WeightedRoundRobinStrategy{},
			"WEIGHTED_RANDOM":                &algo_default.WeightedRandomStrategy{},
			"LOW_LATENCY_WEIGHTED_BALANCER":  &custom.LowLatencyWeightedBalancerStrategy{},
			"WEIGHTED_SUCCESS_RATE_BALANCER": &custom.WeightedSuccessRateBalancer{},
		},
	}
}

func (f *LoadBalancerFactory) GetStrategy(name string) LoadBalancerStrategy {
	return f.strategies[name]
}

var Factory LoadBalancerFactory
