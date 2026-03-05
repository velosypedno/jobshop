package naive

import (
	"errors"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
)

var ErrNoChildrenFound = errors.New("no children found")

type OperationSolution struct {
	Operation      *base.Operation
	MachineID      base.MachineID
	Period         base.Period
	ChildSolutions []*OperationSolution
}

func (os *OperationSolution) GetLastChildCompletionTime() (time.Time, error) {
	if len(os.ChildSolutions) == 0 {
		return time.Time{}, ErrNoChildrenFound
	}

	maxTime := os.ChildSolutions[0].Period.End
	for _, child := range os.ChildSolutions {
		if child.Period.End.After(maxTime) {
			maxTime = child.Period.End
		}
	}
	return maxTime, nil
}

type JobSolution struct {
	Job                *base.Job
	OperationSolutions []*OperationSolution
}

type Solution struct {
	Jobs []*JobSolution
}

func (s *Solution) ToBaseSolution() *base.Solution {
	baseJobs := make([]*base.JobSolution, len(s.Jobs))
	for i, js := range s.Jobs {
		baseJobs[i] = js.toBaseJobSolution()
	}

	return &base.Solution{
		Jobs: baseJobs,
	}
}

func (js *JobSolution) toBaseJobSolution() *base.JobSolution {
	baseOpSolutions := make([]*base.OperationSolution, len(js.OperationSolutions))
	for i, os := range js.OperationSolutions {
		baseOpSolutions[i] = os.toBaseOperationSolution()
	}

	return &base.JobSolution{
		Job:                js.Job,
		OperationSolutions: baseOpSolutions,
	}
}

func (os *OperationSolution) toBaseOperationSolution() *base.OperationSolution {
	baseChildSolutions := make([]*base.OperationSolution, len(os.ChildSolutions))
	for i, child := range os.ChildSolutions {
		baseChildSolutions[i] = child.toBaseOperationSolution()
	}

	return &base.OperationSolution{
		Operation:      os.Operation,
		MachineID:      os.MachineID,
		Period:         os.Period,
		ChildSolutions: baseChildSolutions,
	}
}
