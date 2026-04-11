package base

import "time"

type Problem struct {
	Jobs      []*Job
	Machines  []*Machine
	StartTime time.Time
}

type OperationSolutionV2 struct {
	MachineID MachineID
	Offset    time.Duration
	StartTime time.Time
}

type SolutionV2 struct {
	OperationMap map[OperationID]OperationSolutionV2
}

func NewSolutionV2() SolutionV2 {
	return SolutionV2{
		OperationMap: make(map[OperationID]OperationSolutionV2),
	}
}
