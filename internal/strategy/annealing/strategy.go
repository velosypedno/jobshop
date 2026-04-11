package annealing

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/simulator"
	"go.uber.org/zap"
)

type Strategy struct {
	InitialTemp      float64
	MinTemp          float64
	Alpha            float64
	Iterations       int
	SwapsPerMutation int
	logger           *zap.Logger
	name             string
}

func New(initialTemp, minTemp, alpha float64, iterations int, swaps int, name string) *Strategy {
	l, _ := zap.NewProduction()
	return &Strategy{
		InitialTemp:      initialTemp,
		MinTemp:          minTemp,
		Alpha:            alpha,
		Iterations:       iterations,
		SwapsPerMutation: swaps,
		logger:           l,
		name:             name,
	}
}

func (s *Strategy) SetLogger(l *zap.Logger) {
	s.logger = l
}

func (s *Strategy) Type() string {
	return "Simulated Annealing (Priority-Based)"
}

func (s *Strategy) Name() string {
	return s.name
}

func (s *Strategy) Description() string {
	return fmt.Sprintf(
		"Optimization using thermodynamic annealing with Priority Weighting.\n"+
			"It evolves a vector of weights for each operation. The simulator picks the\n"+
			"best candidate from the ReadyList based on these priorities.\n\n"+
			"| %-18s | %-10s |\n"+
			"|:-------------------|-----------:|\n"+
			"| %-18s | %10.2f |\n"+
			"| %-18s | %10.4f |\n"+
			"| %-18s | %10.4f |\n"+
			"| %-18s | %10d |\n"+
			"| %-18s | %10d |",
		"Parameter", "Value",
		"Initial Temp", s.InitialTemp,
		"Min Temp", s.MinTemp,
		"Alpha (Cooling)", s.Alpha,
		"Iterations / T", s.Iterations,
		"Swaps Per Mutate", s.SwapsPerMutation,
	)
}

func (s *Strategy) Plan(problem base.Problem) (*base.Solution, base.MachineTimeSlots) {
	sim := simulator.NewFactorySimulator(problem)
	n := sim.TotalOperations()

	s.logger.Info("Starting Simulated Annealing",
		zap.String("strategy_type", s.Type()),
		zap.Float64("initial_temp", s.InitialTemp),
		zap.Float64("alpha", s.Alpha),
		zap.Int("ops_count", n),
	)

	currentWeights := make([]float64, n)
	for i := range currentWeights {
		currentWeights[i] = rand.Float64()
	}

	currentRes := sim.Simulate(currentWeights)
	bestRes := currentRes
	temp := s.InitialTemp

	for temp > s.MinTemp {
		s.logger.Info("Temperature cycle",
			zap.Float64("temp", temp),
			zap.Any("current_cost", currentRes.Cost),
			zap.Any("best_cost", bestRes.Cost),
		)

		for i := 0; i < s.Iterations; i++ {
			nextWeights := s.mutate(currentWeights)
			nextRes := sim.Simulate(nextWeights)

			if s.shouldAccept(float64(currentRes.Cost), float64(nextRes.Cost), temp) {
				currentRes = nextRes
				if currentRes.Cost < bestRes.Cost {
					bestRes = currentRes
					s.logger.Debug("Global best updated",
						zap.Any("cost", bestRes.Cost),
						zap.Float64("at_temp", temp),
					)
				}
			}
		}
		temp *= s.Alpha
	}

	s.logger.Info("Simulated Annealing finished",
		zap.Any("final_cost", bestRes.Cost),
		zap.Duration("elapsed", time.Since(problem.StartTime)),
	)

	return bestRes.Solution, bestRes.MachineSlots
}

func (s *Strategy) mutate(weights []float64) []float64 {
	next := make([]float64, len(weights))
	copy(next, weights)

	for i := 0; i < s.SwapsPerMutation; i++ {
		idx1 := rand.Intn(len(next))
		idx2 := rand.Intn(len(next))
		next[idx1], next[idx2] = next[idx2], next[idx1]
	}

	return next
}

func (s *Strategy) shouldAccept(current, next, temp float64) bool {
	if next < current {
		return true
	}
	probability := math.Exp((current - next) / temp)
	return rand.Float64() < probability
}
