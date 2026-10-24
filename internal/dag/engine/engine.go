package engine

import (
	"errors"
	"time"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

type Report struct {
	Solution        *core.Solution
	StrategyMetrics StrategyMetrics
	SolutionMetrics SolutionMetrics
}

type StrategyMetrics struct {
	StrategyType        string
	StrategyName        string
	StrategyDescription string
	SchedulingTime      time.Duration
}

type SolutionMetrics struct {
	MakeSpan         time.Duration
	UtilizationLevel float64
}

var ErrNoStrategySet = errors.New("no strategy set")

type Engine struct {
	strategy core.Strategy
}

func New(strategy core.Strategy) *Engine {
	return &Engine{
		strategy: strategy,
	}
}

func (e *Engine) Solve(p *core.Problem) (Report, error) {
	if e.strategy == nil {
		return Report{}, ErrNoStrategySet
	}

	startPlanning := time.Now()
	solution := e.strategy.Plan(p)

	metaInfo := StrategyMetrics{
		StrategyName:        e.strategy.Name(),
		StrategyType:        e.strategy.Type(),
		StrategyDescription: e.strategy.Description(),
		SchedulingTime:      time.Since(startPlanning),
	}

	workflowPeriod := solution.GetPeriod()
	makeSpan := workflowPeriod.Duration
	utilization := 0.0
	if makeSpan > 0 {
		utilization = solution.GerUtilizationLevel()
	}

	report := Report{
		Solution: &solution,
		SolutionMetrics: SolutionMetrics{
			MakeSpan:         makeSpan,
			UtilizationLevel: utilization,
		},
		StrategyMetrics: metaInfo,
	}

	return report, nil
}
