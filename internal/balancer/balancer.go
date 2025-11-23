package balancer

import (
	"load-balancer/internal/backend"
)

type LoadBalancingAlgorithm interface {
	GetNextBackend(pool *backend.ServerPool) *backend.Backend
	Name() string
}

type RoundRobbinAlgorithm struct {}

func (r *RoundRobbinAlgorithm) GetNextBackend(pool *backend.ServerPool) *backend.Backend {
	return pool.GetNextBackend()
}

func (r *RoundRobbinAlgorithm) Name() string {
	return "round-robin"
}

type AlgorithmFactory struct {}

func (f *AlgorithmFactory) GetAlgorithm(name string) LoadBalancingAlgorithm {
	switch name {
	case "round-robin":
		return &RoundRobbinAlgorithm{}
	default:
		return nil
	}
}