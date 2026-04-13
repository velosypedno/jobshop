package simulator

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
