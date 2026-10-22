package core

import (
	"fmt"
)

// readonly
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

// allow machines with different names, but equal types
type MachineConfigEntry struct {
	Type  MachineType
	Name  string
	Count int
}
