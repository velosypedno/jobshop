package tabu

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/simulator"
)

type Strategy struct {
	TabuSize       int
	MaxIterations  int
	NeighborsCount int
}

func New(tabuSize, maxIter, neighbors int) *Strategy {
	return &Strategy{
		TabuSize:       tabuSize,
		MaxIterations:  maxIter,
		NeighborsCount: neighbors,
	}
}

func (s *Strategy) Name() string {
	return "Tabu Search (Priority-Based)"
}

func (s *Strategy) Description() string {
	return fmt.Sprintf(
		"Optimization using Tabu Search with short-term memory.\n"+
			"It prevents the algorithm from returning to recently visited states.\n\n"+
			"| %-18s | %-10s |\n"+
			"|:-------------------|-----------:|\n"+
			"| %-18s | %10d |\n"+
			"| %-18s | %10d |\n"+
			"| %-18s | %10d |",
		"Parameter", "Value",
		"Tabu List Size", s.TabuSize,
		"Max Iterations", s.MaxIterations,
		"Neighbors/Step", s.NeighborsCount,
	)
}

type move struct {
	i, j int
}

func (s *Strategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (*base.Solution, base.MachineTimeSlots) {
	sim := simulator.NewFactorySimulator(jobs, machines, startTime)
	n := sim.TotalOperations()

	currentWeights := make([]float64, n)
	for i := range currentWeights {
		currentWeights[i] = rand.Float64()
	}

	currentRes := sim.Simulate(currentWeights)
	bestRes := currentRes

	tabuList := make(map[move]int)

	for iter := 0; iter < s.MaxIterations; iter++ {
		var bestNeighborRes *simulator.SimulationResult
		var bestNeighborWeights []float64
		var chosenMove move

		for k := 0; k < s.NeighborsCount; k++ {
			i, j := rand.Intn(n), rand.Intn(n)
			m := move{i, j}

			candidateWeights := make([]float64, n)
			copy(candidateWeights, currentWeights)
			candidateWeights[i], candidateWeights[j] = candidateWeights[j], candidateWeights[i]

			res := sim.Simulate(candidateWeights)

			isTabu := tabuList[m] > iter
			if !isTabu || res.Cost < bestRes.Cost {
				if bestNeighborRes == nil || res.Cost < bestNeighborRes.Cost {
					bestNeighborRes = res
					bestNeighborWeights = candidateWeights
					chosenMove = m
				}
			}
		}

		if bestNeighborRes != nil {
			currentWeights = bestNeighborWeights
			currentRes = bestNeighborRes

			tabuList[chosenMove] = iter + s.TabuSize

			if currentRes.Cost < bestRes.Cost {
				bestRes = currentRes
			}
		}
	}

	return bestRes.Solution, bestRes.MachineSlots
}
