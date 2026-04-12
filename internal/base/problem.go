package base

import "time"

type Problem struct {
	Jobs      []*Job
	Machines  []*Machine
	StartTime time.Time
}
