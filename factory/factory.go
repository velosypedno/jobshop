package factory

import (
	"fmt"
	"time"

	"github.com/velosypedno/resource-allocation/base"
	"github.com/velosypedno/resource-allocation/parser"
)

type PlannerStrategy interface {
	Plan([]*base.Job, []*base.Machine, time.Time) (base.Solution, base.MachineTimeSlots)
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

func (f *Factory) AddJobByName(templateName string) error {
	template, ok := f.Templates[templateName]
	if !ok {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	f.jobCounter++
	newJob := base.CreateJob(base.JobID(f.jobCounter), template)
	f.Jobs = append(f.Jobs, &newJob)
	fmt.Println(f.Jobs)
	return nil
}

func (f *Factory) SetPlanner(planner PlannerStrategy) {
	f.Planner = planner
}

func (f *Factory) Plan(startTime time.Time) (base.Solution, SchedulingInfo, error) {
	if f.Planner == nil {
		return base.Solution{}, SchedulingInfo{}, fmt.Errorf("planner strategy is not set")
	}

	startPlanning := time.Now()
	fmt.Println(f.Jobs)
	solution, machineSlotsMap := f.Planner.Plan(f.Jobs, f.Machines, startTime)

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
