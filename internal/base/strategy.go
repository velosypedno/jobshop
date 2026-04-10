package base

import (
	"time"

	"go.uber.org/zap"
)

type Strategy interface {
	Plan([]*Job, []*Machine, time.Time) (*Solution, MachineTimeSlots)
	Type() string
	Name() string
	Description() string
	SetLogger(l *zap.Logger)
}
