package naive

import (
	"errors"
	"time"

	"github.com/velosypedno/resource-allocation/base"
)

const name = "Greedy"
const description = `Greedy Earliest Completion Time scheduling. Each operation is assigned to the machine that 
provides the earliest completion time, taking into account the technological sequence 
(dependence on child operations) and already occupied time slots.`

type Strategy struct{}

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
	session := newSession(machines, startTime)

	solution := Solution{}
	for _, job := range jobs {
		jobSolution := planJob(job, session)
		solution.Jobs = append(solution.Jobs, jobSolution)
	}

	return solution.ToBaseSolution(), session.OccupiedMap
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
