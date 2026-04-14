package engine

import (
	"errors"
	"time"

	"github.com/velosypedno/jobshop/internal/core"
	"go.uber.org/zap"
)

var ErrNoStrategiesSet = errors.New("no strategies set")

type Engine struct {
	strategies []core.Strategy
}

func New(strategies ...core.Strategy) *Engine {
	return &Engine{
		strategies: strategies,
	}
}

func (s *Engine) SetLogger(l *zap.Logger) {
	for _, s := range s.strategies {
		s.SetLogger(l)
	}

}

func (e *Engine) Solve(p *core.Problem) ([]Report, error) {
	if len(e.strategies) == 0 {
		return nil, ErrNoStrategiesSet
	}

	results := make([]Report, 0, len(e.strategies))

	for _, planner := range e.strategies {
		startPlanning := time.Now()
		solution := planner.Plan(p)

		metaInfo := StrategyMetrics{
			StrategyName:        planner.Name(),
			StrategyType:        planner.Type(),
			StrategyDescription: planner.Description(),
			SchedulingTime:      time.Since(startPlanning),
		}

		workflowPeriod := solution.GetPeriod(p.StartTime)
		makeSpan := workflowPeriod.Duration()
		utilization := 0.0
		if makeSpan > 0 {
			utilization = solution.GerUtilizationLevel(p.StartTime)
		}

		results = append(results, Report{
			Solution: &solution,
			SolutionMetrics: SolutionMetrics{
				MakeSpan:         makeSpan,
				UtilizationLevel: utilization,
			},
			StrategyMetrics: metaInfo,
		})
	}

	return results, nil
}

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
