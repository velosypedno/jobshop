package simulator

import (
	"fmt"
	"time"

	"github.com/velosypedno/jobshop/pkg/tree/core"
)

type SimulationResult struct {
	Cost     float64
	Solution core.Solution
}

// A FactorySimulator converts problem with tree-based jobs into a shcedule
// using priory-based algorithm
type FactorySimulator struct {
	// ops stores all flattened problem operations
	// The index of this slice corresponds strictly to internalOp.ID
	ops     []*internalOp
	problem *core.Problem // to access machines, jobs and start time
}

type internalOp struct {
	ID          core.OperationID
	BaseOp      *core.Operation
	JobID       core.JobID
	ParentID    core.OperationID
	InDegree    int
	ChildrenIDs []core.OperationID
}

func (o internalOp) String() string {
	parentInfo := "NONE"
	if o.ParentID != -1 {
		parentInfo = fmt.Sprintf("%d", o.ParentID)
	}

	return fmt.Sprintf(
		"[Op %-3d | Job %-3d] %-15s | Type: %v | InDegree: %d | Parent: %-4s | Children: %v",
		o.ID,
		o.JobID,
		o.BaseOp.Name,
		o.BaseOp.MachineType,
		o.InDegree,
		parentInfo,
		o.ChildrenIDs,
	)
}

// Be sure that problem will not be changed, so you do not need to do copy
func NewFactorySimulator(problem *core.Problem) *FactorySimulator {
	sim := &FactorySimulator{
		ops:     []*internalOp{},
		problem: problem,
	}
	sim.flattenJobs(problem.Jobs)
	return sim
}

// Use to identidy how long should be slice of weigths of operations for simulation
func (s *FactorySimulator) TotalOperations() int {
	return len(s.ops)
}

func (s *FactorySimulator) flattenJobs(jobs []*core.Job) {
	// I have no idea why I needed so strange map
	registry := make(map[core.JobID]map[core.OperationID]*internalOp)

	globalIDCounter := core.OperationID(0)

	var registerRecursive func(jobID core.JobID, ops []*core.Operation)
	registerRecursive = func(jobID core.JobID, ops []*core.Operation) {
		if _, ok := registry[jobID]; !ok {
			registry[jobID] = make(map[core.OperationID]*internalOp)
		}

		for _, op := range ops {
			internal := &internalOp{
				ID:          globalIDCounter,
				BaseOp:      op,
				JobID:       jobID,
				ParentID:    -1,
				InDegree:    len(op.ChildOperations),
				ChildrenIDs: make([]core.OperationID, 0, len(op.ChildOperations)),
			}

			s.ops = append(s.ops, internal)
			registry[jobID][op.ID] = internal
			globalIDCounter++

			registerRecursive(jobID, op.ChildOperations)
		}
	}

	for _, job := range jobs {
		registerRecursive(job.ID, job.Operations)
	}

	for _, job := range jobs {
		var linkRecursive func(ops []*core.Operation)
		linkRecursive = func(ops []*core.Operation) {
			for _, parentOp := range ops {
				parentInternal := registry[job.ID][parentOp.ID]

				for _, childOp := range parentOp.ChildOperations {
					childInternal := registry[job.ID][childOp.ID]

					childInternal.ParentID = parentInternal.ID
					parentInternal.ChildrenIDs = append(parentInternal.ChildrenIDs, childInternal.ID)

					linkRecursive([]*core.Operation{childOp})
				}
			}
		}
		linkRecursive(job.Operations)
	}
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
		currentInDegrees[i] = op.InDegree
	}

	readyList := make([]core.OperationID, 0, total)

	for i, deg := range currentInDegrees {
		if deg == 0 {
			readyList = append(readyList, core.OperationID(i))
		}
	}

	sess := newSession(s.problem.Machines)
	var maxFinishOffset time.Duration = 0

	for len(readyList) > 0 {
		bestPos := pickBestOperation(readyList, weights)
		opIdx := readyList[bestPos]
		op := s.ops[opIdx]

		readyList[bestPos] = readyList[len(readyList)-1]
		readyList = readyList[:len(readyList)-1]

		opSolution := sess.FindOpSulotion(op)

		if opSolution.Period.End() > maxFinishOffset {
			maxFinishOffset = opSolution.Period.End()
		}

		if op.ParentID != -1 {
			currentInDegrees[op.ParentID]--
			if currentInDegrees[op.ParentID] == 0 {
				readyList = append(readyList, op.ParentID)
			}
		}

	}
	solution := s.assemble(sess)
	return &SimulationResult{
		Cost:     maxFinishOffset.Seconds(),
		Solution: solution,
	}
}

func pickBestOperation(readyList []core.OperationID, weights []float64) int {
	bestIdx := 0
	for i := 1; i < len(readyList); i++ {
		if weights[readyList[i]] < weights[readyList[bestIdx]] {
			bestIdx = i
		}
	}
	return bestIdx
}

func (s *FactorySimulator) assemble(sess *session) core.Solution {
	solution := core.NewSolution()

	for _, op := range s.ops {
		opSolution, ok := sess.opSolutions[op.ID]

		if !ok {
			panic(fmt.Sprintf("No solution found in session for operation with ID: %v", op.ID))
		}

		solution.OperationMap[op.BaseOp.ID] = opSolution
	}

	return solution
}
