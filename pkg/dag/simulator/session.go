package simulator

import (
	"sort"
	"time"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

type MachineTimeSlots map[core.MachineID][]core.Period
type MachineTypeIndex map[core.MachineType][]core.MachineID

type session struct {
	occupiedMap      MachineTimeSlots
	machineTypeIndex MachineTypeIndex

	Solution core.Solution
}

func newSession(machines []*core.Machine) *session {
	timeSlotsMap := make(map[core.MachineID][]core.Period)
	for _, machine := range machines {
		timeSlotsMap[machine.ID] = []core.Period{}
	}

	machineTypeIndex := make(map[core.MachineType][]core.MachineID)
	for _, machine := range machines {
		machineTypeIndex[machine.Type] = append(machineTypeIndex[machine.Type], machine.ID)
	}

	solution := core.Solution{
		OperationMap: make(map[core.OpID]core.OpSolution),
	}
	return &session{
		occupiedMap:      timeSlotsMap,
		machineTypeIndex: machineTypeIndex,

		Solution: solution,
	}
}

func (s *session) FindOpSulotion(
	op *core.Operation,
) core.OpSolution {
	startOffset := s.getReadyOffset(op)
	targetMachineIDs := s.machineTypeIndex[op.MachineType]

	var bestMachineID core.MachineID
	var bestPeriod core.Period
	firstFound := false

	for _, mID := range targetMachineIDs {
		currentPeriod := s.findEarliestGap(startOffset, op.ProcessingTime, s.occupiedMap[mID])

		if !firstFound || currentPeriod.End() < bestPeriod.End() {
			bestPeriod = currentPeriod
			bestMachineID = mID
			firstFound = true
		}
	}

	opSolution := core.OpSolution{
		Period:    bestPeriod,
		MachineID: bestMachineID,
	}

	s.Solution.OperationMap[op.ID] = opSolution
	s.occupiedMap[bestMachineID] = append(s.occupiedMap[bestMachineID], bestPeriod)

	return opSolution
}

func (s *session) findEarliestGap(startOffset time.Duration, duration time.Duration, occupied []core.Period) core.Period {
	sort.Slice(occupied, func(i, j int) bool {
		return occupied[i].Offset < occupied[j].Offset
	})

	candidateOffset := startOffset

	for _, slot := range occupied {
		if slot.End() < candidateOffset || slot.End() == candidateOffset {
			continue
		}

		if slot.Offset-candidateOffset >= duration {
			return core.Period{
				Offset:   candidateOffset,
				Duration: duration,
			}
		}

		if slot.End() > candidateOffset {
			candidateOffset = slot.End()
		}
	}

	return core.Period{
		Offset:   candidateOffset,
		Duration: duration,
	}
}

func (s *session) getReadyOffset(op *core.Operation) time.Duration {
	var readyOffset time.Duration = 0

	for _, predecessorID := range op.Predecessors {
		if predecessorOpSoluiton, ok := s.Solution.OperationMap[predecessorID]; ok {
			if predecessorOpSoluiton.Period.End() > readyOffset {
				readyOffset = predecessorOpSoluiton.Period.End()
			}
		}
	}
	return readyOffset
}
