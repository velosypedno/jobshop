package base

import "time"

type Period struct {
	Start time.Time
	End   time.Time
}

func (p Period) Duration() time.Duration {
	return p.End.Sub(p.Start)
}

type OpSolution struct {
	MachineID MachineID
	Offset    time.Duration
	Duration  time.Duration
}

type Solution struct {
	OperationMap map[OperationID]OpSolution
}

func NewSolution() Solution {
	return Solution{
		OperationMap: make(map[OperationID]OpSolution),
	}
}

func (s *Solution) GetPeriod(startTime time.Time) Period {
	var maxEndTime time.Time

	for _, opSol := range s.OperationMap {
		endTime := startTime.Add(opSol.Offset + opSol.Duration)
		if endTime.After(maxEndTime) {
			maxEndTime = endTime
		}
	}

	return Period{
		Start: startTime,
		End:   maxEndTime,
	}
}

func (s *Solution) GetAllOperationsDuration() time.Duration {
	duration := time.Duration(0)
	for _, opSol := range s.OperationMap {
		duration += opSol.Duration
	}
	return duration
}

func (s *Solution) GerUtilizationLevel(startTime time.Time) float64 {
	period := s.GetPeriod(startTime)

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

	return float64(duration) / (float64(machinesCount) * float64(period.Duration()))
}
