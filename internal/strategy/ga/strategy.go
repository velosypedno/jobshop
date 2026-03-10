package ga

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/simulator"
)

type Strategy struct {
	PopulationSize int
	Generations    int
	MutationRate   float64
	CrossoverRate  float64
	ElitismRatio   float64
}

func New(popSize, generations int, mutationRate, crossoverRate, elitism float64) *Strategy {
	return &Strategy{
		PopulationSize: popSize,
		Generations:    generations,
		MutationRate:   mutationRate,
		CrossoverRate:  crossoverRate,
		ElitismRatio:   elitism,
	}
}

func (s *Strategy) Name() string {
	return "Genetic Algorithm (Priority-Based)"
}

func (s *Strategy) Description() string {
	return fmt.Sprintf(
		"Optimization using Evolutionary Computing with Priority Weighting.\n"+
			"It evolves a population of weight vectors using selection, crossover,\n"+
			"and mutation to find the most efficient production sequence.\n\n"+
			"| %-18s | %-10s |\n"+
			"|:-------------------|-----------:|\n"+
			"| %-18s | %10d |\n"+
			"| %-18s | %10d |\n"+
			"| %-18s | %10.2f |\n"+
			"| %-18s | %10.2f |\n"+
			"| %-18s | %10.2f |",
		"Parameter", "Value",
		"Population Size", s.PopulationSize,
		"Generations", s.Generations,
		"Mutation Rate", s.MutationRate,
		"Crossover Rate", s.CrossoverRate,
		"Elitism Ratio", s.ElitismRatio,
	)
}

type individual struct {
	weights []float64
	fitness float64
	result  *simulator.SimulationResult
}

func (s *Strategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (*base.Solution, base.MachineTimeSlots) {
	sim := simulator.NewFactorySimulator(jobs, machines, startTime)
	n := sim.TotalOperations()

	population := make([]*individual, s.PopulationSize)
	for i := 0; i < s.PopulationSize; i++ {
		weights := make([]float64, n)
		for j := 0; j < n; j++ {
			weights[j] = rand.Float64()
		}

		res := sim.Simulate(weights)
		population[i] = &individual{
			weights: weights,
			result:  res,
			fitness: 1.0 / float64(res.Cost+1),
		}
	}

	for g := 0; g < s.Generations; g++ {
		sort.Slice(population, func(i, j int) bool {
			return population[i].fitness > population[j].fitness
		})

		if g%10 == 0 {
			fmt.Printf("[%s] Generation %d: Best Makespan = %v\n",
				s.Name(), g, population[0].result.Cost)
		}

		newPopulation := make([]*individual, 0, s.PopulationSize)

		elitismCount := int(float64(s.PopulationSize) * s.ElitismRatio)
		if elitismCount < 1 {
			elitismCount = 1
		}
		newPopulation = append(newPopulation, population[:elitismCount]...)

		for len(newPopulation) < s.PopulationSize {
			p1 := s.selectParent(population)
			p2 := s.selectParent(population)

			childWeights := s.crossover(p1, p2)

			s.mutate(childWeights)

			res := sim.Simulate(childWeights)
			newPopulation = append(newPopulation, &individual{
				weights: childWeights,
				result:  res,
				fitness: 1.0 / float64(res.Cost+1),
			})
		}

		population = newPopulation
	}

	sort.Slice(population, func(i, j int) bool {
		return population[i].fitness > population[j].fitness
	})

	best := population[0]
	fmt.Printf("[%s] Final Best: %v\n", s.Name(), best.result.Cost)

	return best.result.Solution, best.result.MachineSlots
}
func (s *Strategy) selectParent(population []*individual) *individual {
	tournamentSize := 3
	var best *individual

	for i := 0; i < tournamentSize; i++ {
		randomIndex := rand.Intn(len(population))
		contender := population[randomIndex]

		if best == nil || contender.fitness > best.fitness {
			best = contender
		}
	}
	return best
}

func (s *Strategy) crossover(p1, p2 *individual) []float64 {
	n := len(p1.weights)
	childWeights := make([]float64, n)

	if rand.Float64() < s.CrossoverRate {
		pivot := rand.Intn(n)

		for i := 0; i < n; i++ {
			if i < pivot {
				childWeights[i] = p1.weights[i]
			} else {
				childWeights[i] = p2.weights[i]
			}
		}
	} else {
		if p1.fitness > p2.fitness {
			copy(childWeights, p1.weights)
		} else {
			copy(childWeights, p2.weights)
		}
	}

	return childWeights
}

func (s *Strategy) mutate(weights []float64) {
	for i := range weights {
		if rand.Float64() < s.MutationRate {
			weights[i] += rand.NormFloat64() * 0.1

			if weights[i] < 0 {
				weights[i] = 0
			}
			if weights[i] > 1 {
				weights[i] = 1
			}
		}
	}
}
