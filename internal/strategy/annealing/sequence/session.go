package sequence

import (
	"sort"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
)

type fullID struct {
	jobID base.JobID
	opID  base.OperationID
}

type session struct {
	OccupiedMap      base.MachineTimeSlots
	MachineTypeIndex base.MachineTypeIndex
	StartTime        time.Time

	results map[fullID]base.Period
}

func newSession(machines []*base.Machine, startTime time.Time) *session {
	return &session{
		OccupiedMap:      initTimeSlotsMap(machines),
		MachineTypeIndex: initMachineTypeIndex(machines),
		StartTime:        startTime,
		results:          make(map[fullID]base.Period, 0),
	}
}

func initTimeSlotsMap(machines []*base.Machine) base.MachineTimeSlots {
	timeSlotsMap := make(map[base.MachineID][]base.Period)
	for _, machine := range machines {
		timeSlotsMap[machine.ID] = []base.Period{}
	}
	return timeSlotsMap
}

func initMachineTypeIndex(machines []*base.Machine) base.MachineTypeIndex {
	machineTypeIndex := make(map[base.MachineType][]base.MachineID)
	for _, machine := range machines {
		machineTypeIndex[machine.Type] = append(machineTypeIndex[machine.Type], machine.ID)
	}
	return machineTypeIndex
}

func (s *session) FindBestSlot(
	startTime time.Time,
	duration time.Duration,
	machineType base.MachineType,
) (base.MachineID, base.Period) {
	targetMachineIDs := s.MachineTypeIndex[machineType]

	var bestMachineID base.MachineID
	var bestPeriod base.Period
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

func (s *session) findEarliestGap(startTime time.Time, duration time.Duration, occupied []base.Period) base.Period {
	sort.Slice(occupied, func(i, j int) bool {
		return occupied[i].Start.Before(occupied[j].Start)
	})

	candidateStart := startTime

	for _, slot := range occupied {
		if slot.End.Before(candidateStart) {
			continue
		}
		if slot.Start.Sub(candidateStart) >= duration {
			return base.Period{
				Start: candidateStart,
				End:   candidateStart.Add(duration),
			}
		}

		if slot.End.After(candidateStart) {
			candidateStart = slot.End
		}
	}

	return base.Period{
		Start: candidateStart,
		End:   candidateStart.Add(duration),
	}
}

func (s *session) GetReadyTime(op *base.Operation) time.Time {
	readyTime := s.StartTime
	for _, child := range op.ChildOperations {
		key := fullID{jobID: op.JobID, opID: child.ID}
		if childPeriod, ok := s.results[key]; ok {
			if childPeriod.End.After(readyTime) {
				readyTime = childPeriod.End
			}
		}
	}
	return readyTime
}
