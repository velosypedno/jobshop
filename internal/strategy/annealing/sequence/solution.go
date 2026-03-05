package sequence

import "github.com/velosypedno/resource-allocation/internal/base"

type OperationSolution struct {
	Operation      *base.Operation
	MachineID      base.MachineID
	Period         base.Period
	ChildSolutions []*OperationSolution
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
