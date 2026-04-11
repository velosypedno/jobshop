package base

import (
	"go.uber.org/zap"
)

type Strategy interface {
	Plan(Problem) (*Solution, MachineTimeSlots)
	Type() string
	Name() string
	Description() string
	SetLogger(l *zap.Logger)
}
