package naive

import (
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"go.uber.org/zap"
)

const strategyType = "Greedy"
const description = `Greedy Earliest Completion Time scheduling. Each operation is assigned to the machine that 
provides the earliest completion time, taking into account the technological sequence 
(dependence on child operations) and already occupied time slots.`

type Strategy struct {
	logger *zap.Logger
	name   string
}

func New(name string) *Strategy {
	l, _ := zap.NewProduction()
	return &Strategy{
		logger: l,
		name:   name,
	}
}

func (s *Strategy) SetLogger(l *zap.Logger) {
	s.logger = l
}
func (Strategy) Type() string {
	return strategyType
}

func (s *Strategy) Name() string {
	return s.name
}

func (Strategy) Description() string {
	return description
}

func (s *Strategy) Plan(problem *base.Problem) base.SolutionV2 {
	s.logger.Info("Starting Greedy planning",
		zap.String("strategy_type", s.Type()),
		zap.Int("jobs_count", len(problem.Jobs)),
		zap.Int("machines_count", len(problem.Machines)),
	)

	session := newSession(problem.Machines, problem.StartTime)

	solV2 := base.NewSolutionV2()

	for _, job := range problem.Jobs {
		planJob(job, session, &solV2)
	}

	s.logger.Info("Greedy planning completed",
		zap.Duration("elapsed", time.Since(problem.StartTime)),
	)

	return solV2
}

func planJob(
	job *base.Job,
	session *session,
	solution *base.SolutionV2,
) {
	for _, operation := range job.Operations {
		planOperation(operation, session, solution)
	}
}

func planOperation(
	operation *base.Operation,
	session *session,
	solution *base.SolutionV2,
) {
	operationSolutionV2 := base.OperationSolutionV2{}

	for _, child := range operation.ChildOperations {
		planOperation(child, session, solution)
	}

	lastChildEndTime := session.StartTime
	for _, childOp := range operation.ChildOperations {
		if childOp == nil {
			continue
		}
		childSolution := solution.OperationMap[childOp.ID]
		childCompletionTime := session.StartTime.Add(childSolution.Offset).Add(childOp.Duration)
		if childCompletionTime.After(lastChildEndTime) {
			lastChildEndTime = childCompletionTime
		}
	}

	targetMachineID, targetPeriod := session.FindBestSlot(
		lastChildEndTime,
		operation.Duration,
		operation.MachineType,
	)

	operationSolutionV2.Duration = operation.Duration
	operationSolutionV2.MachineID = targetMachineID
	operationSolutionV2.Offset = targetPeriod.Start.Sub(session.StartTime)
	solution.OperationMap[operation.ID] = operationSolutionV2

	session.OccupiedMap[targetMachineID] = append(session.OccupiedMap[targetMachineID], targetPeriod)

}
