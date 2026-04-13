package naive

import (
	"errors"
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
	solution := Solution{}

	for _, job := range problem.Jobs {
		jobSolution := planJob(job, session)
		solution.Jobs = append(solution.Jobs, jobSolution)
	}

	s.logger.Info("Greedy planning completed",
		zap.Duration("elapsed", time.Since(problem.StartTime)),
	)

	solV2 := base.NewSolutionV2()

	for _, jobSolution := range solution.Jobs {
		for _, operationSolution := range jobSolution.toBaseJobSolution().GetAllOperations() {
			solV2.OperationMap[operationSolution.Operation.ID] = base.OperationSolutionV2{
				MachineID: operationSolution.MachineID,
				Offset:    operationSolution.Period.Start.Sub(problem.StartTime),
				Duration:  operationSolution.Operation.Duration,
			}
		}
	}

	return solV2
}

func planJob(
	job *base.Job,
	session *session,
) *JobSolution {
	jobSolution := &JobSolution{
		Job:                job,
		OperationSolutions: []*OperationSolution{},
	}
	for _, operation := range job.Operations {
		operationSolution := planOperation(operation, session)
		jobSolution.OperationSolutions = append(jobSolution.OperationSolutions, operationSolution)
	}
	return jobSolution
}

func planOperation(
	operation *base.Operation,
	session *session,
) *OperationSolution {
	operationSolution := &OperationSolution{
		Operation:      operation,
		ChildSolutions: []*OperationSolution{},
	}

	for _, child := range operation.ChildOperations {
		operationSolution.ChildSolutions = append(
			operationSolution.ChildSolutions,
			planOperation(child, session))
	}
	lastChildEndTime, err := operationSolution.GetLastChildCompletionTime()
	if errors.Is(err, ErrNoChildrenFound) {
		targetMachineID, targetPeriod := session.FindBestSlot(
			session.StartTime,
			operation.Duration,
			operation.MachineType,
		)
		operationSolution.MachineID = targetMachineID
		operationSolution.Period = targetPeriod
		session.OccupiedMap[targetMachineID] = append(session.OccupiedMap[targetMachineID], targetPeriod)
		return operationSolution
	}
	if err != nil {
		panic(err)
	}

	targetMachineID, targetPeriod := session.FindBestSlot(
		lastChildEndTime,
		operation.Duration,
		operation.MachineType,
	)
	operationSolution.MachineID = targetMachineID
	operationSolution.Period = targetPeriod
	session.OccupiedMap[targetMachineID] = append(session.OccupiedMap[targetMachineID], targetPeriod)

	return operationSolution
}
