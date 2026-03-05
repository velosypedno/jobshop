package rnd

import (
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
)

const name = "Random Search"
const description = `Generates a random sequence of operations (Sequence) 
and schedules them based on the Earliest Slot principle, strictly adhering 
to technological dependencies (post-order traversal).`

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func (Strategy) Name() string        { return name }
func (Strategy) Description() string { return description }

func (s *Strategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (*base.Solution, base.MachineTimeSlots) {
	sess := newSession(machines, startTime)

	counts := make([]int, len(jobs))
	for i, job := range jobs {
		counts[i] = job.OperationsCount()
	}
	seq := NewSequence(counts)
	seq.Shuffle()

	jobCounters := make([]int, len(jobs))
	plannedOps := make(map[fullID]*OperationSolution)

	for i := 0; i < seq.Len(); i++ {
		jobIdx := seq.Get(i)
		job := jobs[jobIdx]
		opIdx := jobCounters[jobIdx]
		operation := job.GetOperation(opIdx)

		readyTime := sess.GetReadyTime(operation)
		mID, period := sess.FindBestSlot(readyTime, operation.Duration, operation.MachineType)

		currentID := fullID{jobID: operation.JobID, opID: operation.ID}

		sess.results[currentID] = period

		opSol := &OperationSolution{
			Operation:      operation,
			MachineID:      mID,
			Period:         period,
			ChildSolutions: []*OperationSolution{},
		}

		for _, child := range operation.ChildOperations {
			childID := fullID{jobID: operation.JobID, opID: child.ID}
			if childSol, ok := plannedOps[childID]; ok {
				opSol.ChildSolutions = append(opSol.ChildSolutions, childSol)
			}
		}

		plannedOps[currentID] = opSol
		sess.OccupiedMap[mID] = append(sess.OccupiedMap[mID], period)
		jobCounters[jobIdx]++
	}

	return s.assemble(jobs, plannedOps), sess.OccupiedMap
}

func (s *Strategy) assemble(
	jobs []*base.Job,
	plannedOps map[fullID]*OperationSolution,
) *base.Solution {
	localSolution := Solution{
		Jobs: make([]*JobSolution, 0, len(jobs)),
	}

	for _, job := range jobs {
		js := &JobSolution{
			Job:                job,
			OperationSolutions: []*OperationSolution{},
		}

		for _, rootOp := range job.Operations {
			key := fullID{jobID: job.ID, opID: rootOp.ID}
			if sol, ok := plannedOps[key]; ok {
				js.OperationSolutions = append(js.OperationSolutions, sol)
			}
		}
		localSolution.Jobs = append(localSolution.Jobs, js)
	}

	return localSolution.ToBaseSolution()
}
