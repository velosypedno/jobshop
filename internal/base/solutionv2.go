package base

import "time"

type OperationSolutionV2 struct {
	MachineID MachineID
	Offset    time.Duration
	Duration  time.Duration
}

type SolutionV2 struct {
	OperationMap map[OperationID]OperationSolutionV2
}

func NewSolutionV2() SolutionV2 {
	return SolutionV2{
		OperationMap: make(map[OperationID]OperationSolutionV2),
	}
}

func (s *SolutionV2) GetPeriod(startTime time.Time) Period {
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

func (s *SolutionV2) GetAllOperationsDuration() time.Duration {
	duration := time.Duration(0)
	for _, opSol := range s.OperationMap {
		duration += opSol.Duration
	}
	return duration
}

func (s *SolutionV2) GerUtilizationLevel(startTime time.Time) float64 {
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
