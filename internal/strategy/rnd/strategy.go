package rnd

import (
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"go.uber.org/zap"
)

const strategyType = "Random Search"
const description = `Generates a random sequence of operations (Sequence) 
and schedules them based on the Earliest Slot principle, strictly adhering 
to technological dependencies (post-order traversal).`

type Strategy struct {
	name string
}

func (_ Strategy) SetLogger(_ *zap.Logger) {}

func New(name string) *Strategy {
	return &Strategy{
		name: name,
	}
}
func (s *Strategy) Name() string     { return s.name }
func (Strategy) Type() string        { return strategyType }
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
	plannedOps := make(map[base.OperationID]*OperationSolution)

	for i := 0; i < seq.Len(); i++ {
		jobIdx := seq.Get(i)
		job := jobs[jobIdx]
		opIdx := jobCounters[jobIdx]
		operation := job.GetOperation(opIdx)

		readyTime := sess.GetReadyTime(operation)
		mID, period := sess.FindBestSlot(readyTime, operation.Duration, operation.MachineType)

		sess.results[operation.ID] = period

		opSol := &OperationSolution{
			Operation:      operation,
			MachineID:      mID,
			Period:         period,
			ChildSolutions: []*OperationSolution{},
		}

		for _, child := range operation.ChildOperations {
			if childSol, ok := plannedOps[child.ID]; ok {
				opSol.ChildSolutions = append(opSol.ChildSolutions, childSol)
			}
		}

		plannedOps[operation.ID] = opSol
		sess.OccupiedMap[mID] = append(sess.OccupiedMap[mID], period)
		jobCounters[jobIdx]++
	}

	return s.assemble(jobs, plannedOps), sess.OccupiedMap
}

func (s *Strategy) assemble(
	jobs []*base.Job,
	plannedOps map[base.OperationID]*OperationSolution,
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
			if sol, ok := plannedOps[rootOp.ID]; ok {
				js.OperationSolutions = append(js.OperationSolutions, sol)
			}
		}
		localSolution.Jobs = append(localSolution.Jobs, js)
	}

	return localSolution.ToBaseSolution()
}
