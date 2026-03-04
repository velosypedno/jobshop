package factory

import (
	"fmt"
	"time"

	"github.com/velosypedno/resource-allocation/base"
	"github.com/velosypedno/resource-allocation/parser"
)

type PlannerStrategy interface {
	Plan([]*base.Job, []*base.Machine, time.Time) (*base.Solution, base.MachineTimeSlots)
	Name() string
	Description() string
}

type Factory struct {
	Jobs      []*base.Job
	Machines  []*base.Machine
	Templates map[string]base.JobTemplate
	Planner   PlannerStrategy

	machineTypeRegistry map[string]base.MachineType

	jobCounter     int
	machineCounter int
}

func (f *Factory) Configure(machineConfigs []parser.MachineConfig, templates []base.JobTemplate) {
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

func (f *Factory) SetPlanner(planner PlannerStrategy) {
	f.Planner = planner
}

func (f *Factory) Plan(orders []parser.OrderDTO, startTime time.Time) (*base.Solution, SchedulingInfo, error) {
	if f.Planner == nil {
		return &base.Solution{}, SchedulingInfo{}, fmt.Errorf("planner strategy is not set")
	}

	startPlanning := time.Now()

	jobs, err := f.createJobsFromOrders(orders)
	if err != nil {
		return &base.Solution{}, SchedulingInfo{}, err
	}
	solution, machineSlotsMap := f.Planner.Plan(jobs, f.Machines, startTime)

	metaInfo := SchedulingMetaInfo{
		StrategyName:        f.Planner.Name(),
		StrategyDescription: f.Planner.Description(),
		SchedulingTime:      time.Since(startPlanning),
	}

	workflowPeriod := solution.GetWorkFlowPeriod()
	makeSpan := workflowPeriod.Duration()

	utilization := 0.0
	if makeSpan > 0 {
		utilization = machineSlotsMap.GetUtilizationLevel(makeSpan)
	}

	return solution, SchedulingInfo{
		SchedulingMetaInfo: metaInfo,
		MakeSpan:           makeSpan,
		UtilizationLevel:   utilization,
	}, nil
}

func (f *Factory) createJobsFromOrders(orders []parser.OrderDTO) ([]*base.Job, error) {
	var jobs []*base.Job
	jobIDCounter := 0

	for _, order := range orders {
		template, ok := f.Templates[order.Name]
		if !ok {
			return nil, fmt.Errorf("template '%s' not found for order", order.Name)
		}

		for i := 0; i < order.Amount; i++ {
			jobIDCounter++
			newJob := base.CreateJob(base.JobID(jobIDCounter), template)
			jobs = append(jobs, &newJob)
		}
	}
	return jobs, nil
}

type SchedulingMetaInfo struct {
	StrategyName        string
	StrategyDescription string
	SchedulingTime      time.Duration
}

type SchedulingInfo struct {
	SchedulingMetaInfo
	MakeSpan         time.Duration
	UtilizationLevel float64
}
