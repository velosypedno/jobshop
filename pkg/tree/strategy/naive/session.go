package naive

import (
	"sort"
	"time"

	"github.com/velosypedno/jobshop/pkg/tree/core"
)

type session struct {
	OccupiedMap      core.MachineTimeSlots
	MachineTypeIndex core.MachineTypeIndex
}

func newSession(machines []*core.Machine) *session {
	return &session{
		OccupiedMap:      initTimeSlotsMap(machines),
		MachineTypeIndex: initMachineTypeIndex(machines),
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

func (s *session) findEarliestGap(startOffest time.Duration, duration time.Duration, occupied []core.Period) core.Period {
	sort.Slice(occupied, func(i, j int) bool {
		return occupied[i].Offset < occupied[j].Offset
	})

	candidateOffset := startOffest

	for _, slot := range occupied {
		if slot.End() < candidateOffset {
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
