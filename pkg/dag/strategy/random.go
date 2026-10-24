package strategy

import (
	"math/rand/v2"

	"github.com/velosypedno/jobshop/pkg/dag/core"
	"github.com/velosypedno/jobshop/pkg/dag/simulator"
)

type RandomStrategy struct{}

func (s *RandomStrategy) Type() string {
	return "Random planner (Priority-Based)"
}

func (s *RandomStrategy) Name() string {
	return "Random planner (Priority-Based)"
}

func (s *RandomStrategy) Description() string {
	return "Assigns random priority weights to operations to generate a schedule."
}

func (s *RandomStrategy) Plan(problem *core.Problem) core.Solution {
	sim := simulator.NewFactorySimulator(problem)
	n := sim.TotalOperations()

	currentWeights := make([]float64, n)
	for i := range currentWeights {
		currentWeights[i] = rand.Float64()
	}

	currentRes := sim.Simulate(currentWeights)
	return currentRes.Solution
}
