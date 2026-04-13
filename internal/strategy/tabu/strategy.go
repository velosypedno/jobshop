package tabu

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/simulator"
	"go.uber.org/zap"
)

type Strategy struct {
	TabuSize       int
	MaxIterations  int
	NeighborsCount int
	logger         *zap.Logger
	name           string
}

func New(tabuSize, maxIter, neighbors int, name string) *Strategy {
	l, _ := zap.NewProduction()
	return &Strategy{
		TabuSize:       tabuSize,
		MaxIterations:  maxIter,
		NeighborsCount: neighbors,
		logger:         l,
		name:           name,
	}
}

func (s *Strategy) SetLogger(l *zap.Logger) {
	s.logger = l
}

func (s *Strategy) Type() string {
	return "Tabu Search (Priority-Based)"
}

func (s *Strategy) Name() string {
	return s.name
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

func (s *Strategy) Plan(problem *base.Problem) base.SolutionV2 {
	sim := simulator.NewFactorySimulator(problem)
	n := sim.TotalOperations()

	s.logger.Info("Starting Tabu Search optimization",
		zap.String("strategy_type", s.Type()),
		zap.Int("max_iterations", s.MaxIterations),
		zap.Int("neighbors_count", s.NeighborsCount),
		zap.Int("tabu_size", s.TabuSize),
	)

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
				s.logger.Debug("New global best found",
					zap.Int("iteration", iter),
					zap.Any("cost", bestRes.Cost),
				)
			}
		}

		s.logger.Info("Iteration completed",
			zap.Int("iter", iter),
			zap.Any("current_cost", currentRes.Cost),
			zap.Any("best_cost", bestRes.Cost),
		)
	}

	s.logger.Info("Tabu Search finished",
		zap.Any("final_best_cost", bestRes.Cost),
		zap.Duration("elapsed", time.Since(problem.StartTime)),
	)

	return bestRes.SolutionV2
}
