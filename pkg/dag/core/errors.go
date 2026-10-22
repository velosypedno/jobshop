package core

import (
	"errors"
	"fmt"
)

var ErrCycleInGraph = errors.New("cycle in graph")
var ErrUnknownOpTemplate = errors.New("unknown op template")
var ErrUnknownJobTemplate = errors.New("unknown job template")
var ErrMissingMachineType = errors.New("missing machine type")

type ErrUnknownOpTemplateDetail struct {
	ID OpTemplateID
}

func (e *ErrUnknownOpTemplateDetail) Error() string {
	return fmt.Sprintf("op template with id %v does not exist", e.ID)
}

// Unwrap allows errors.Is(err, ErrUnknownOpTemplate) to work.
func (e *ErrUnknownOpTemplateDetail) Unwrap() error { return ErrUnknownOpTemplate }

type ErrMissingMachineTypeDetail struct {
	MachineType MachineType
}

func (e *ErrMissingMachineTypeDetail) Error() string {
	return fmt.Sprintf("machine type %v is not configured in the manufactory", e.MachineType)
}

// Unwrap allows errors.Is(err, ErrMissingMachineType) to work.
func (e *ErrMissingMachineTypeDetail) Unwrap() error { return ErrMissingMachineType }

type ErrUnknownJobTemplateDetail struct {
	ID JobTemplateID
}

func (e *ErrUnknownJobTemplateDetail) Error() string {
	return fmt.Sprintf("job template with id %v not found", e.ID)
}

// Unwrap allows errors.Is(err, ErrUnknownJobTemplate) to work.
func (e *ErrUnknownJobTemplateDetail) Unwrap() error { return ErrUnknownJobTemplate }

func IsCycleInGraph(err error) bool {
	return errors.Is(err, ErrCycleInGraph)
}

func IsMissingMachineType(err error) bool {
	return errors.Is(err, ErrMissingMachineType)
}

func IsUnknownOpTemplate(err error) bool {
	return errors.Is(err, ErrUnknownOpTemplate)
}

func IsUnknownJobTemplate(err error) bool {
	return errors.Is(err, ErrUnknownJobTemplate)
}

func AsUnknownOpTemplate(err error) (*ErrUnknownOpTemplateDetail, bool) {
	var detail *ErrUnknownOpTemplateDetail
	ok := errors.As(err, &detail)
	return detail, ok
}

func AsMissingMachineType(err error) (*ErrMissingMachineTypeDetail, bool) {
	var detail *ErrMissingMachineTypeDetail
	ok := errors.As(err, &detail)
	return detail, ok
}

func AsUnknownJobTemplate(err error) (*ErrUnknownJobTemplateDetail, bool) {
	var detail *ErrUnknownJobTemplateDetail
	ok := errors.As(err, &detail)
	return detail, ok
}
