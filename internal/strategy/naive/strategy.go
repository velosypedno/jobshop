package naive

import (
	"errors"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"go.uber.org/zap"
)

const name = "Greedy"
const description = `Greedy Earliest Completion Time scheduling. Each operation is assigned to the machine that 
provides the earliest completion time, taking into account the technological sequence 
(dependence on child operations) and already occupied time slots.`

type Strategy struct {
	logger *zap.Logger
}

func New() *Strategy {
	l, _ := zap.NewProduction()
	return &Strategy{
		logger: l,
	}
}

func (s *Strategy) SetLogger(l *zap.Logger) {
	s.logger = l
}
func (Strategy) Name() string {
	return name
}

func (Strategy) Description() string {
	return description
}

func (s *Strategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (*base.Solution, base.MachineTimeSlots) {
	s.logger.Info("Starting Greedy planning",
		zap.String("strategy_type", s.Name()),
		zap.Int("jobs_count", len(jobs)),
		zap.Int("machines_count", len(machines)),
	)

	session := newSession(machines, startTime)
	solution := Solution{}

	for _, job := range jobs {
		jobSolution := planJob(job, session)
		solution.Jobs = append(solution.Jobs, jobSolution)
	}

	baseSolution := solution.ToBaseSolution()

	s.logger.Info("Greedy planning completed",
		zap.Duration("elapsed", time.Since(startTime)),
	)

	return baseSolution, session.OccupiedMap
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
