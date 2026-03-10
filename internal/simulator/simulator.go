package simulator

import (
	"fmt"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
)

type SimulationResult struct {
	Cost         float64
	Solution     *base.Solution
	MachineSlots base.MachineTimeSlots
}

type FactorySimulator struct {
	Ops          []*InternalOp
	machines     []*base.Machine
	startTime    time.Time
	rootOpIDs    map[base.JobID][]int
	originalJobs []*base.Job
}

type InternalOp struct {
	ID          int
	BaseOp      *base.Operation
	JobID       int
	ParentID    int
	InDegree    int
	ChildrenIDs []int
	MachineID   string
}

func (o InternalOp) String() string {
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

func NewFactorySimulator(jobs []*base.Job, machines []*base.Machine, startTime time.Time) *FactorySimulator {
	sim := &FactorySimulator{
		Ops:          []*InternalOp{},
		machines:     machines,
		startTime:    startTime,
		rootOpIDs:    make(map[base.JobID][]int),
		originalJobs: jobs,
	}
	sim.flattenJobs(jobs)
	return sim
}

func (s *FactorySimulator) TotalOperations() int {
	return len(s.Ops)
}
func (s *FactorySimulator) flattenJobs(jobs []*base.Job) {
	registry := make(map[base.JobID]map[base.OperationID]*InternalOp)

	globalIDCounter := 0

	var registerRecursive func(jobID base.JobID, ops []*base.Operation)
	registerRecursive = func(jobID base.JobID, ops []*base.Operation) {
		if _, ok := registry[jobID]; !ok {
			registry[jobID] = make(map[base.OperationID]*InternalOp)
		}

		for _, op := range ops {
			internal := &InternalOp{
				ID:          globalIDCounter,
				BaseOp:      op,
				JobID:       int(jobID),
				ParentID:    -1,
				InDegree:    len(op.ChildOperations),
				ChildrenIDs: make([]int, 0, len(op.ChildOperations)),
			}

			s.Ops = append(s.Ops, internal)
			registry[jobID][op.ID] = internal
			globalIDCounter++

			registerRecursive(jobID, op.ChildOperations)
		}
	}

	for _, job := range jobs {
		registerRecursive(job.ID, job.Operations)
		for _, rootOp := range job.Operations {
			if internal, ok := registry[job.ID][rootOp.ID]; ok {
				s.rootOpIDs[job.ID] = append(s.rootOpIDs[job.ID], internal.ID)
			}
		}
	}

	for _, job := range jobs {
		var linkRecursive func(ops []*base.Operation)
		linkRecursive = func(ops []*base.Operation) {
			for _, parentOp := range ops {
				parentInternal := registry[job.ID][parentOp.ID]

				for _, childOp := range parentOp.ChildOperations {
					childInternal := registry[job.ID][childOp.ID]

					childInternal.ParentID = parentInternal.ID
					parentInternal.ChildrenIDs = append(parentInternal.ChildrenIDs, childInternal.ID)

					linkRecursive([]*base.Operation{childOp})
				}
			}
		}
		linkRecursive(job.Operations)
	}
}

func (s *FactorySimulator) Simulate(weights []float64) *SimulationResult {
	total := len(s.Ops)
	if total == 0 {
		return &SimulationResult{Solution: &base.Solution{}}
	}

	currentInDegrees := make([]int, total)
	for i, op := range s.Ops {
		currentInDegrees[i] = op.InDegree
	}

	readyList := make([]int, 0, total)

	for i, deg := range currentInDegrees {
		if deg == 0 {
			readyList = append(readyList, i)
		}
	}

	sess := newSession(s.machines, s.startTime)
	var maxFinishTime time.Time = s.startTime

	for len(readyList) > 0 {
		bestPos := pickBestOperation(readyList, weights)
		opIdx := readyList[bestPos]
		op := s.Ops[opIdx]

		readyList[bestPos] = readyList[len(readyList)-1]
		readyList = readyList[:len(readyList)-1]

		readyTime := sess.GetReadyTime(op)

		mID, period := sess.FindBestSlot(readyTime, op.BaseOp.Duration, op.BaseOp.MachineType)

		sess.results[op.ID] = period
		sess.assignedMachines[op.ID] = mID
		sess.OccupiedMap[mID] = append(sess.OccupiedMap[mID], period)

		if period.End.After(maxFinishTime) {
			maxFinishTime = period.End
		}

		if op.ParentID != -1 {
			currentInDegrees[op.ParentID]--
			if currentInDegrees[op.ParentID] == 0 {
				readyList = append(readyList, op.ParentID)
			}
		}

	}

	return &SimulationResult{
		Cost:         maxFinishTime.Sub(s.startTime).Seconds(),
		MachineSlots: sess.OccupiedMap,
		Solution:     s.Assemble(sess),
	}
}

func pickBestOperation(readyList []int, weights []float64) int {
	bestIdx := 0
	for i := 1; i < len(readyList); i++ {
		if weights[readyList[i]] < weights[readyList[bestIdx]] {
			bestIdx = i
		}
	}
	return bestIdx
}

func (s *FactorySimulator) Assemble(sess *session) *base.Solution {
	opSols := make(map[int]*OperationSolution, len(s.Ops))

	for _, op := range s.Ops {
		period, ok := sess.results[op.ID]
		mID, okM := sess.assignedMachines[op.ID]

		if !ok || !okM {
			continue
		}

		opSols[op.ID] = &OperationSolution{
			Operation:      op.BaseOp,
			MachineID:      mID,
			Period:         period,
			ChildSolutions: []*OperationSolution{},
		}
	}

	for _, op := range s.Ops {
		if op.ParentID != -1 {
			parentSol := opSols[op.ParentID]
			childSol := opSols[op.ID]
			if parentSol != nil && childSol != nil {
				parentSol.ChildSolutions = append(parentSol.ChildSolutions, childSol)
			}
		}
	}

	localSolution := &Solution{
		Jobs: make([]*JobSolution, 0, len(s.originalJobs)),
	}

	for _, job := range s.originalJobs {
		js := &JobSolution{
			Job:                job,
			OperationSolutions: []*OperationSolution{},
		}

		for _, rootID := range s.rootOpIDs[job.ID] {
			if sol, ok := opSols[rootID]; ok {
				js.OperationSolutions = append(js.OperationSolutions, sol)
			}
		}
		localSolution.Jobs = append(localSolution.Jobs, js)
	}

	return localSolution.ToBaseSolution()
}
