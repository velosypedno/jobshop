package simulator

import (
	"sort"
	"time"

	"github.com/velosypedno/jobshop/pkg/tree/core"
)

type session struct {
	OccupiedMap      core.MachineTimeSlots
	MachineTypeIndex core.MachineTypeIndex

	results          map[core.OperationID]core.Period
	assignedMachines map[core.OperationID]core.MachineID
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

	return &session{
		OccupiedMap:      timeSlotsMap,
		MachineTypeIndex: machineTypeIndex,
		results:          make(map[core.OperationID]core.Period, 0),
		assignedMachines: make(map[core.OperationID]core.MachineID),
	}
}

func (s *session) FindBestSlot(
	startOffset time.Duration,
	duration time.Duration,
	machineType core.MachineType,
) (core.MachineID, core.Period) {
	targetMachineIDs := s.MachineTypeIndex[machineType]

	var bestMachineID core.MachineID
	var bestPeriod core.Period
	firstFound := false

	for _, mID := range targetMachineIDs {
		currentPeriod := s.findEarliestGap(startOffset, duration, s.OccupiedMap[mID])

		if !firstFound || currentPeriod.End() < bestPeriod.End() {
			bestPeriod = currentPeriod
			bestMachineID = mID
			firstFound = true
		}
	}

	return bestMachineID, bestPeriod
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

func (s *session) GetReadyOffset(op *internalOp) time.Duration {
	var readyOffset time.Duration = 0

	for _, childGlobalID := range op.ChildrenIDs {
		if childPeriod, ok := s.results[childGlobalID]; ok {
			if childPeriod.End() > readyOffset {
				readyOffset = childPeriod.End()
			}
		}
	}
	return readyOffset
}
