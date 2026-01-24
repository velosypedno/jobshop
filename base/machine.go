package base

import (
	"fmt"
	"time"
)

type (
	MachineType uint
	MachineID   int
)

type MachineTimeSlots map[MachineID][]Period
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
