package scheduler

import (
	"fmt"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/parser"
	"go.uber.org/zap"
)

type PlanResult struct {
	Solution *base.Solution
	Info     SchedulingInfo
}

type Scheduler struct {
	Jobs      []*base.Job
	Machines  []*base.Machine
	Templates map[string]base.JobTemplate

	Planners []base.Strategy

	machineTypeRegistry map[string]base.MachineType
	jobCounter          int
	operationCounter    base.OperationID
	machineCounter      int
}

func (s *Scheduler) SetLogger(l *zap.Logger) {

	for _, planner := range s.Planners {
		planner.SetLogger(l)
	}

}

func (f *Scheduler) Configure(machineConfigs []parser.MachineConfig, templates []base.JobTemplate) {
	f.Templates = make(map[string]base.JobTemplate)
	f.machineTypeRegistry = make(map[string]base.MachineType)

	for _, t := range templates {
		f.Templates[t.Name] = t
	}

	for _, mConf := range machineConfigs {
		mType := base.MachineType(mConf.TypeID)
		f.machineTypeRegistry[mConf.TypeName] = mType

		for i := 0; i < mConf.Count; i++ {
			f.machineCounter++
			m := base.NewMachine(base.MachineID(f.machineCounter), mType, mConf.TypeName)
			m.Name = mConf.TypeName
			f.Machines = append(f.Machines, &m)
		}
	}
}

func (f *Scheduler) SetPlanners(planners ...base.Strategy) {
	f.Planners = planners
}

func (f *Scheduler) Plan(orders []parser.OrderDTO, startTime time.Time) ([]PlanResult, error) {
	if len(f.Planners) == 0 {
		return nil, fmt.Errorf("no planner strategies set")
	}

	jobs, err := f.createJobsFromOrders(orders)
	if err != nil {
		return nil, err
	}

	problem := base.Problem{
		Jobs:      jobs,
		Machines:  f.Machines,
		StartTime: startTime,
	}

	results := make([]PlanResult, 0, len(f.Planners))

	for _, planner := range f.Planners {
		startPlanning := time.Now()

		solution, solutionV2 := planner.Plan(problem)

		metaInfo := SchedulingMetaInfo{
			StrategyName:        planner.Name(),
			StrategyType:        planner.Type(),
			StrategyDescription: planner.Description(),
			SchedulingTime:      time.Since(startPlanning),
		}

		workflowPeriod := solution.GetWorkFlowPeriod()
		makeSpan := workflowPeriod.Duration()
		utilization := 0.0
		if makeSpan > 0 {
			utilization = solutionV2.GerUtilizationLevel(startTime)
		}

		results = append(results, PlanResult{
			Solution: solution,
			Info: SchedulingInfo{
				SchedulingMetaInfo: metaInfo,
				MakeSpan:           makeSpan,
				UtilizationLevel:   utilization,
			},
		})
	}

	return results, nil
}

func (f *Scheduler) createJobsFromOrders(orders []parser.OrderDTO) ([]*base.Job, error) {
	var jobs []*base.Job
	jobIDCounter := 0

	for _, order := range orders {
		template, ok := f.Templates[order.Name]
		if !ok {
			return nil, fmt.Errorf("template '%s' not found for order", order.Name)
		}

		for i := 0; i < order.Amount; i++ {
			jobIDCounter++
			newJob := base.CreateJob(base.JobID(jobIDCounter), template, &f.operationCounter)
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
