package sequence

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"go.uber.org/zap"
)

type Strategy struct {
	InitialTemp      float64
	MinTemp          float64
	Alpha            float64
	Iterations       int
	SwapsPerMutation int
}

func (_ Strategy) SetLogger(_ *zap.Logger) {}

func New(initialTemp, minTemp, alpha float64, iterations int, swaps int) *Strategy {
	return &Strategy{
		InitialTemp:      initialTemp,
		MinTemp:          minTemp,
		Alpha:            alpha,
		Iterations:       iterations,
		SwapsPerMutation: swaps,
	}
}

func (s *Strategy) Name() string {
	return "Simulated Annealing (Sequence-Based)"
}

func (s *Strategy) Description() string {
	return fmt.Sprintf(
		"Optimization using thermodynamic annealing with Sequence Encoding.\n"+
			"It evolves a fixed global queue of operations (e.g., 001021), where the simulator\n"+
			"strictly follows the order of jobs in the sequence.\n\n"+
			"| %-18s | %-10s |\n"+
			"|:-------------------|-----------:|\n"+
			"| %-18s | %10.2f |\n"+
			"| %-18s | %10.4f |\n"+
			"| %-18s | %10.4f |\n"+
			"| %-18s | %10d |\n"+
			"| %-18s | %10d |",
		"Parameter", "Value",
		"Initial Temp", s.InitialTemp,
		"Min Temp", s.MinTemp,
		"Alpha (Cooling)", s.Alpha,
		"Iterations / T", s.Iterations,
		"Swaps Per Mutate", s.SwapsPerMutation,
	)
}
func (s *Strategy) calculateCost(sol *base.Solution) float64 {
	period := sol.GetWorkFlowPeriod()
	return period.End.Sub(period.Start).Seconds()
}

func (s *Strategy) shouldAccept(currentCost, nextCost, temp float64) bool {
	if nextCost <= currentCost {
		return true
	}
	delta := nextCost - currentCost
	return rand.Float64() < math.Exp(-delta/temp)
}

func (s *Strategy) Plan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
) (*base.Solution, base.MachineTimeSlots) {
	counts := make([]int, len(jobs))
	for i, job := range jobs {
		counts[i] = job.OperationsCount()
	}
	currentSeq := NewSequence(counts)
	currentSeq.Shuffle()

	currentSol, currentSlots := s.runInternalPlan(jobs, machines, startTime, currentSeq)
	currentCost := s.calculateCost(currentSol)

	bestSol := currentSol
	bestSlots := currentSlots
	bestCost := currentCost

	temp := s.InitialTemp
	for temp > s.MinTemp {
		for i := 0; i < s.Iterations; i++ {
			nextSeq := s.mutate(currentSeq)

			nextSol, nextSlots := s.runInternalPlan(jobs, machines, startTime, nextSeq)
			nextCost := s.calculateCost(nextSol)

			if s.shouldAccept(currentCost, nextCost, temp) {
				currentSol = nextSol
				currentSlots = nextSlots
				currentCost = nextCost

				if currentCost < bestCost {
					bestSol = currentSol
					bestSlots = currentSlots
					bestCost = currentCost
				}
			}
		}

		temp *= s.Alpha
	}

	return bestSol, bestSlots
}

func (s *Strategy) runInternalPlan(
	jobs []*base.Job,
	machines []*base.Machine,
	startTime time.Time,
	seq *Sequence,
) (*base.Solution, base.MachineTimeSlots) {
	sess := newSession(machines, startTime)
	jobCounters := make([]int, len(jobs))
	plannedOps := make(map[fullID]*OperationSolution)

	for i := 0; i < seq.Len(); i++ {
		jobIdx := seq.Get(i)
		job := jobs[jobIdx]
		opIdx := jobCounters[jobIdx]
		operation := job.GetOperation(opIdx)

		readyTime := sess.GetReadyTime(operation)
		mID, period := sess.FindBestSlot(readyTime, operation.Duration, operation.MachineType)

		curID := fullID{jobID: operation.JobID, opID: operation.ID}
		sess.results[curID] = period

		opSol := &OperationSolution{
			Operation: operation,
			MachineID: mID,
			Period:    period,
		}

		for _, child := range operation.ChildOperations {
			childID := fullID{jobID: operation.JobID, opID: child.ID}
			if childSol, ok := plannedOps[childID]; ok {
				opSol.ChildSolutions = append(opSol.ChildSolutions, childSol)
			}
		}

		plannedOps[curID] = opSol
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

func (s *Strategy) mutate(original *Sequence) *Sequence {
	next := original.Clone()
	n := next.Len()

	if n < 2 {
		return next
	}

	successfulSwaps := 0
	for attempt := 0; attempt < s.SwapsPerMutation*10 && successfulSwaps < s.SwapsPerMutation; attempt++ {
		i := rand.Intn(n)
		j := rand.Intn(n)

		if i != j && next.Get(i) != next.Get(j) {
			next.Swap(i, j)
			successfulSwaps++
		}
	}

	return next
}
