package priority

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/simulator"
)

type Simulator interface {
	Simulate(weights []float64) *simulator.SimulationResult
	TotalOperations() int
}

type Strategy struct {
	InitialTemp      float64
	MinTemp          float64
	Alpha            float64
	Iterations       int
	SwapsPerMutation int
}

func New(initialTemp, minTemp, alpha float64, iterations int, swaps int) *Strategy {
	return &Strategy{
		InitialTemp:      initialTemp,
		MinTemp:          minTemp,
		Alpha:            alpha,
		Iterations:       iterations,
		SwapsPerMutation: swaps,
	}
}

func (s *Strategy) Name() string {
	return "Simulated Annealing (Priority-Based)"
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
	temp := s.InitialTemp
	for temp > s.MinTemp {
		for i := 0; i < s.Iterations; i++ {
			nextWeights := s.mutate(currentWeights)

			nextRes := sim.Simulate(nextWeights)

			if s.shouldAccept(currentRes.Cost, nextRes.Cost, temp) {
				currentRes = nextRes

				if currentRes.Cost < bestRes.Cost {
					bestRes = currentRes
				}
			}
		}
		temp *= s.Alpha
	}

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
