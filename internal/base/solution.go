package base

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrNoChildrenFound = errors.New("no children found")

type Period struct {
	Start time.Time
	End   time.Time
}

func (p Period) Duration() time.Duration {
	return p.End.Sub(p.Start)
}

type OperationSolution struct {
	Operation      *Operation
	MachineID      MachineID
	Period         Period
	ChildSolutions []*OperationSolution
}

func (os *OperationSolution) Flatten() []*OperationSolution {
	results := []*OperationSolution{os}
	for _, child := range os.ChildSolutions {
		results = append(results, child.Flatten()...)
	}
	return results
}

func (os *OperationSolution) GetSelfDuration() time.Duration {
	return os.Operation.Duration
}

func (os *OperationSolution) GetTreeWorkDuration() time.Duration {
	total := os.GetSelfDuration()
	for _, child := range os.ChildSolutions {
		total += child.GetTreeWorkDuration()
	}
	return total
}

func (os *OperationSolution) GetTreeFlowPeriod() Period {
	minStart := os.Period.Start
	maxEnd := os.Period.End

	for _, child := range os.ChildSolutions {
		childPeriod := child.GetTreeFlowPeriod()
		if childPeriod.Start.Before(minStart) {
			minStart = childPeriod.Start
		}
		if childPeriod.End.After(maxEnd) {
			maxEnd = childPeriod.End
		}
	}

	return Period{Start: minStart, End: maxEnd}
}

func (os *OperationSolution) GetLastChildCompletionTime() (time.Time, error) {
	if len(os.ChildSolutions) == 0 {
		return time.Time{}, ErrNoChildrenFound
	}

	maxTime := os.ChildSolutions[0].Period.End
	for _, child := range os.ChildSolutions {
		if child.Period.End.After(maxTime) {
			maxTime = child.Period.End
		}
	}
	return maxTime, nil
}

type JobSolution struct {
	Job                *Job
	OperationSolutions []*OperationSolution
}

func (js *JobSolution) GetTotalWorkDuration() time.Duration {
	var total time.Duration
	for _, opSol := range js.OperationSolutions {
		total += opSol.GetTreeWorkDuration()
	}
	return total
}

func (js *JobSolution) GetAllOperations() []*OperationSolution {
	var allOps []*OperationSolution
	for _, rootOp := range js.OperationSolutions {
		allOps = append(allOps, rootOp.Flatten()...)
	}
	return allOps
}

func (js *JobSolution) GetJobFlowPeriod() Period {
	if len(js.OperationSolutions) == 0 {
		return Period{}
	}

	start := js.OperationSolutions[0].Period.Start
	end := js.OperationSolutions[0].Period.End

	for _, opSol := range js.OperationSolutions {
		p := opSol.GetTreeFlowPeriod()
		if p.Start.Before(start) {
			start = p.Start
		}
		if p.End.After(end) {
			end = p.End
		}
	}
	return Period{Start: start, End: end}
}

type Solution struct {
	Jobs []*JobSolution
}

func (s *Solution) GetWorkDuration() time.Duration {
	var duration time.Duration
	for _, jobSolution := range s.Jobs {
		duration += jobSolution.GetTotalWorkDuration()
	}
	return duration
}

func (s *Solution) GetWorkFlowPeriod() Period {
	if len(s.Jobs) == 0 {
		return Period{}
	}

	firstJobPeriod := s.Jobs[0].GetJobFlowPeriod()
	start := firstJobPeriod.Start
	end := firstJobPeriod.End

	for _, jobSolution := range s.Jobs {
		jobPeriod := jobSolution.GetJobFlowPeriod()

		if jobPeriod.Start.Before(start) {
			start = jobPeriod.Start
		}
		if jobPeriod.End.After(end) {
			end = jobPeriod.End
		}
	}

	return Period{
		Start: start,
		End:   end,
	}
}

func (s Solution) String() string {
	var sb strings.Builder
	totalPeriod := s.GetWorkFlowPeriod()

	sb.WriteString("================================================================================\n")
	sb.WriteString("FACTORY PLAN SUMMARY\n")
	sb.WriteString(fmt.Sprintf("Total Duration: %v\n", s.GetWorkDuration()))
	sb.WriteString(fmt.Sprintf("Flow Period:    %s -> %s\n",
		totalPeriod.Start.Format("15:04:05"),
		totalPeriod.End.Format("15:04:05")))
	sb.WriteString("================================================================================\n\n")

	for _, js := range s.Jobs {
		jobPeriod := js.GetJobFlowPeriod()
		sb.WriteString(fmt.Sprintf("JOB: %s [ID: %v]\n", js.Job.Name, js.Job.ID))
		sb.WriteString(fmt.Sprintf("Period: %s - %s\n",
			jobPeriod.Start.Format("15:04:05"),
			jobPeriod.End.Format("15:04:05")))

		for _, opSol := range js.OperationSolutions {
			sb.WriteString(opSol.formatSolutionTree(1))
		}
		sb.WriteString("--------------------------------------------------------------------------------\n")
	}

	return sb.String()
}

func (os *OperationSolution) formatSolutionTree(level int) string {
	var sb strings.Builder

	var indent string
	if level > 1 {
		indent = strings.Repeat("  │ ", level-1) + "  ├─ "
	} else {
		indent = " ├─ "
	}

	sb.WriteString(fmt.Sprintf("%s[%-5v] %-5s (%-10v) | %s -> %s\n",
		indent,
		os.Operation.ID,
		os.Operation.Name,
		os.Operation.MachineType,
		os.Period.Start.Format("15:04:05"),
		os.Period.End.Format("15:04:05"),
	))

	for _, child := range os.ChildSolutions {
		sb.WriteString(child.formatSolutionTree(level + 1))
	}

	return sb.String()
}
