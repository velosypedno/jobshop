package simulator

import (
	"time"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

type SimulationResult struct {
	Cost     float64
	Solution core.Solution
}

type FactorySimulator struct {
	ops     []*core.Operation
	problem *core.Problem // to access machines, jobs and start time
}

// Be sure that problem will not be changed, so you do not need to do copy
func NewFactorySimulator(problem *core.Problem) *FactorySimulator {
	opsAmount := 0
	for _, job := range problem.Jobs {
		opsAmount += len(job.Ops)
	}

	ops := make([]*core.Operation, opsAmount)
	for _, job := range problem.Jobs {
		for _, op := range job.Ops {
			ops[op.ID] = op
		}
	}

	sim := &FactorySimulator{
		ops:     ops,
		problem: problem,
	}
	return sim
}

// Use to identidy how long should be slice of weigths of operations for simulation
func (s *FactorySimulator) TotalOperations() int {
	return len(s.ops)
}

// Simulate constructs a schedule using the provided priority weights for each operation
// Algorithm chose ready operations with the highest priority (LOWER weight values represent HIGHER priority)
// Then add new operaions to ready list if it is possible
//
// Expects len(weights) match with TotalOperations(). If len is not the same it will panic
func (s *FactorySimulator) Simulate(weights []float64) *SimulationResult {
	total := len(s.ops)
	if total == 0 {
		return &SimulationResult{Solution: core.NewSolution()}
	}

	currentInDegrees := make([]int, total)
	for i, op := range s.ops {
		currentInDegrees[i] = len(op.Predecessors)
	}

	readyList := make([]core.OpID, 0, total)

	for index, deg := range currentInDegrees {
		if deg == 0 {
			readyList = append(readyList, core.OpID(index))
		}
	}

	sess := newSession(s.problem.Machines)
	var finishOffset time.Duration = 0

	for len(readyList) > 0 {
		mostPrioritized := 0
		for i := 1; i < len(readyList); i++ {
			if weights[readyList[i]] < weights[readyList[mostPrioritized]] {
				mostPrioritized = i
			}
		}
		opIdx := readyList[mostPrioritized]
		op := s.ops[opIdx]

		readyList[mostPrioritized] = readyList[len(readyList)-1]
		readyList = readyList[:len(readyList)-1]

		opSolution := sess.FindOpSulotion(op)

		if opSolution.Period.End() > finishOffset {
			finishOffset = opSolution.Period.End()
		}

		for _, successorOpID := range op.Successors {
			currentInDegrees[successorOpID]--
			if currentInDegrees[successorOpID] == 0 {
				readyList = append(readyList, successorOpID)
			}
		}

	}

	return &SimulationResult{
		Cost:     finishOffset.Seconds(),
		Solution: sess.Solution,
	}
}
