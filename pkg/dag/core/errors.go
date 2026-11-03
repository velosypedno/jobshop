package core

import (
	"errors"
	"fmt"
)

var ErrCycleInGraph = errors.New("cycle in graph")

type ErrUnknownOpTemplateDetail struct {
	ID OpTemplateID
}

func (e *ErrUnknownOpTemplateDetail) Error() string {
	return fmt.Sprintf("op template with id %v does not exist", e.ID)
}

type ErrMissingMachineTypeDetail struct {
	MachineType MachineType
}

func (e *ErrMissingMachineTypeDetail) Error() string {
	return fmt.Sprintf("machine type %v is not configured in the manufactory", e.MachineType)
}

type ErrUnknownJobTemplateDetail struct {
	ID JobTemplateID
}

func (e *ErrUnknownJobTemplateDetail) Error() string {
	return fmt.Sprintf("job template with id %v not found", e.ID)
}
