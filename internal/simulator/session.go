package simulator

import (
	"sort"
	"time"

	"github.com/velosypedno/resource-allocation/internal/core"
)

type session struct {
	OccupiedMap      core.MachineTimeSlots
	MachineTypeIndex core.MachineTypeIndex
	StartTime        time.Time

	results          map[core.OperationID]core.Period
	assignedMachines map[core.OperationID]core.MachineID
}

func newSession(machines []*core.Machine, startTime time.Time) *session {
	return &session{
		OccupiedMap:      initTimeSlotsMap(machines),
		MachineTypeIndex: initMachineTypeIndex(machines),
		StartTime:        startTime,
		results:          make(map[core.OperationID]core.Period, 0),
		assignedMachines: make(map[core.OperationID]core.MachineID),
	}
}

func initTimeSlotsMap(machines []*core.Machine) core.MachineTimeSlots {
	timeSlotsMap := make(map[core.MachineID][]core.Period)
	for _, machine := range machines {
		timeSlotsMap[machine.ID] = []core.Period{}
	}
	return timeSlotsMap
}

func initMachineTypeIndex(machines []*core.Machine) core.MachineTypeIndex {
	machineTypeIndex := make(map[core.MachineType][]core.MachineID)
	for _, machine := range machines {
		machineTypeIndex[machine.Type] = append(machineTypeIndex[machine.Type], machine.ID)
	}
	return machineTypeIndex
}

func (s *session) FindBestSlot(
	startTime time.Time,
	duration time.Duration,
	machineType core.MachineType,
) (core.MachineID, core.Period) {
	targetMachineIDs := s.MachineTypeIndex[machineType]

	var bestMachineID core.MachineID
	var bestPeriod core.Period
	firstFound := false

	for _, mID := range targetMachineIDs {
		currentPeriod := s.findEarliestGap(startTime, duration, s.OccupiedMap[mID])

		if !firstFound || currentPeriod.End.Before(bestPeriod.End) {
			bestPeriod = currentPeriod
			bestMachineID = mID
			firstFound = true
		}
	}

	return bestMachineID, bestPeriod
}

func (s *session) findEarliestGap(startTime time.Time, duration time.Duration, occupied []core.Period) core.Period {
	sort.Slice(occupied, func(i, j int) bool {
		return occupied[i].Start.Before(occupied[j].Start)
	})

	candidateStart := startTime

	for _, slot := range occupied {
		if slot.End.Before(candidateStart) || slot.End.Equal(candidateStart) {
			continue
		}

		if slot.Start.Sub(candidateStart) >= duration {
			return core.Period{
				Start: candidateStart,
				End:   candidateStart.Add(duration),
			}
		}

		if slot.End.After(candidateStart) {
			candidateStart = slot.End
		}
	}

	return core.Period{
		Start: candidateStart,
		End:   candidateStart.Add(duration),
	}
}

func (s *session) GetReadyTime(op *InternalOp) time.Time {
	readyTime := s.StartTime

	for _, childGlobalID := range op.ChildrenIDs {
		if childPeriod, ok := s.results[childGlobalID]; ok {
			if childPeriod.End.After(readyTime) {
				readyTime = childPeriod.End
			}
		}
	}
	return readyTime
}
