package scheduler

import (
	"fmt"
	"time"

	"github.com/velosypedno/jobshop/internal/core"
	"github.com/velosypedno/jobshop/internal/parser"
	"go.uber.org/zap"
)

type PlanResult struct {
	SolutionV2 *core.Solution
	Info       SchedulingInfo
}

type Scheduler struct {
	Jobs      []*core.Job
	Machines  []*core.Machine
	Templates map[string]core.JobTemplate

	Planners []core.Strategy

	machineTypeRegistry map[string]core.MachineType
	jobCounter          int
	operationCounter    core.OperationID
	machineCounter      int
}

func (s *Scheduler) SetLogger(l *zap.Logger) {

	for _, planner := range s.Planners {
		planner.SetLogger(l)
	}

}

func (f *Scheduler) Configure(machineConfigs []parser.MachineConfig, templates []core.JobTemplate) {
	f.Templates = make(map[string]core.JobTemplate)
	f.machineTypeRegistry = make(map[string]core.MachineType)

	for _, t := range templates {
		f.Templates[t.Name] = t
	}

	for _, mConf := range machineConfigs {
		mType := core.MachineType(mConf.TypeID)
		f.machineTypeRegistry[mConf.TypeName] = mType

		for i := 0; i < mConf.Count; i++ {
			f.machineCounter++
			m := core.NewMachine(core.MachineID(f.machineCounter), mType, mConf.TypeName)
			m.Name = mConf.TypeName
			f.Machines = append(f.Machines, &m)
		}
	}
}

func (f *Scheduler) SetPlanners(planners ...core.Strategy) {
	f.Planners = planners
}

func (f *Scheduler) GetProblem(orders []parser.OrderDTO, startTime time.Time) *core.Problem {
	jobs, err := f.createJobsFromOrders(orders)
	if err != nil {
		return &core.Problem{}
	}

	problem := core.Problem{
		Jobs:      jobs,
		Machines:  f.Machines,
		StartTime: startTime,
	}
	return &problem
}

func (f *Scheduler) Plan(problem *core.Problem) ([]PlanResult, error) {
	if len(f.Planners) == 0 {
		return nil, fmt.Errorf("no planner strategies set")
	}

	results := make([]PlanResult, 0, len(f.Planners))

	for _, planner := range f.Planners {
		startPlanning := time.Now()
		solution := planner.Plan(problem)

		metaInfo := SchedulingMetaInfo{
			StrategyName:        planner.Name(),
			StrategyType:        planner.Type(),
			StrategyDescription: planner.Description(),
			SchedulingTime:      time.Since(startPlanning),
		}

		workflowPeriod := solution.GetPeriod(problem.StartTime)
		makeSpan := workflowPeriod.Duration()
		utilization := 0.0
		if makeSpan > 0 {
			utilization = solution.GerUtilizationLevel(problem.StartTime)
		}

		results = append(results, PlanResult{
			SolutionV2: &solution,
			Info: SchedulingInfo{
				SchedulingMetaInfo: metaInfo,
				MakeSpan:           makeSpan,
				UtilizationLevel:   utilization,
			},
		})
	}

	return results, nil
}

func (f *Scheduler) createJobsFromOrders(orders []parser.OrderDTO) ([]*core.Job, error) {
	var jobs []*core.Job
	jobIDCounter := 0

	for _, order := range orders {
		template, ok := f.Templates[order.Name]
		if !ok {
			return nil, fmt.Errorf("template '%s' not found for order", order.Name)
		}

		for i := 0; i < order.Amount; i++ {
			jobIDCounter++
			newJob := core.CreateJob(core.JobID(jobIDCounter), template, &f.operationCounter)
			jobs = append(jobs, &newJob)
		}
	}
	return jobs, nil
}

type SchedulingMetaInfo struct {
	StrategyType        string
	StrategyName        string
	StrategyDescription string
	SchedulingTime      time.Duration
}

type SchedulingInfo struct {
	SchedulingMetaInfo
	MakeSpan         time.Duration
	UtilizationLevel float64
}
