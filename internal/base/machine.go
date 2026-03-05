package base

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type (
	MachineType uint
	MachineID   int
)

type MachineTimeSlots map[MachineID][]Period

func (m MachineTimeSlots) String() string {
	var out strings.Builder
	out.WriteString("--- Machine Occupancy Status ---\n")

	keys := make([]MachineID, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, mID := range keys {
		slots := m[mID]
		out.WriteString(fmt.Sprintf("Machine [%v]:\n", mID))

		if len(slots) == 0 {
			out.WriteString("  (empty)\n")
			continue
		}

		sortedSlots := make([]Period, len(slots))
		copy(sortedSlots, slots)
		sort.Slice(sortedSlots, func(i, j int) bool {
			return sortedSlots[i].Start.Before(sortedSlots[j].Start)
		})

		for i, p := range sortedSlots {
			overlapWarn := ""
			if i > 0 && p.Start.Before(sortedSlots[i-1].End) {
				overlapWarn = " [!] OVERLAP"
			}
			out.WriteString(fmt.Sprintf("  - %s -> %s%s\n",
				p.Start.Format("15:04:05"),
				p.End.Format("15:04:05"),
				overlapWarn))
		}
	}
	return out.String()
}

type MachineTypeIndex map[MachineType][]MachineID

type Machine struct {
	ID   MachineID
	Type MachineType
	Name string
}

func (m Machine) String() string {
	return fmt.Sprintf("ID: %d, Type: %s", m.ID, m.Name)
}

func NewMachine(id MachineID, machineType MachineType, machineName string) Machine {
	return Machine{
		ID:   id,
		Type: machineType,
		Name: machineName,
	}
}

func (m MachineTimeSlots) GetUtilizationLevel(duration time.Duration) float64 {
	var sumDuration time.Duration

	for _, slots := range m {
		for _, slot := range slots {
			sumDuration += slot.Duration()
		}
	}

	return (float64(sumDuration) / float64(len(m))) / float64(duration)
}
