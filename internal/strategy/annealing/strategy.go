package annealing

import (
	"github.com/velosypedno/resource-allocation/internal/strategy/annealing/priority"
	"github.com/velosypedno/resource-allocation/internal/strategy/annealing/sequence"
)

type Config struct {
	InitialTemp      float64
	MinTemp          float64
	Alpha            float64
	Iterations       int
	SwapsPerMutation int
}

func NewSequenceBased(cfg Config) *sequence.Strategy {
	return sequence.New(
		cfg.InitialTemp,
		cfg.MinTemp,
		cfg.Alpha,
		cfg.Iterations,
		cfg.SwapsPerMutation,
	)
}

func NewPriorityBased(cfg Config) *priority.Strategy {
	return priority.New(
		cfg.InitialTemp,
		cfg.MinTemp,
		cfg.Alpha,
		cfg.Iterations,
		cfg.SwapsPerMutation,
	)
}
