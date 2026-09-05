package core

import "time"

type Period struct {
	Offset   time.Duration
	Duration time.Duration
}

func (p Period) End() time.Duration {
	return p.Offset + p.Duration
}

type OpSolution struct {
	MachineID MachineID
	Period    Period
}

type Solution struct {
	OperationMap map[OperationID]OpSolution
}

func NewSolution() Solution {
	return Solution{
		OperationMap: make(map[OperationID]OpSolution),
	}
}

func (s *Solution) GetPeriod() Period {
	var maxEnd time.Duration

	for _, opSol := range s.OperationMap {
		endTime := opSol.Period.Offset + opSol.Period.Duration
		if endTime > maxEnd {
			maxEnd = endTime
		}
	}

	return Period{
		Offset:   0,
		Duration: maxEnd,
	}
}

func (s *Solution) GetAllOperationsDuration() time.Duration {
	duration := time.Duration(0)
	for _, opSol := range s.OperationMap {
		duration += opSol.Period.Duration
	}
	return duration
}

func (s *Solution) GerUtilizationLevel() float64 {
	period := s.GetPeriod()

	machinesCount := 0

	machinesSet := make(map[MachineID]struct{})
	for _, opSol := range s.OperationMap {
		_, ok := machinesSet[opSol.MachineID]
		if !ok {
			machinesCount++
			machinesSet[opSol.MachineID] = struct{}{}
		}
	}
	duration := s.GetAllOperationsDuration()

	return float64(duration) / (float64(machinesCount) * float64(period.Duration))
}
