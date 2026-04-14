package scheduler

import (
	"fmt"
	"time"

	"github.com/velosypedno/jobshop/internal/core"
	"github.com/velosypedno/jobshop/internal/engine"
	"github.com/velosypedno/jobshop/internal/parser"
	"go.uber.org/zap"
)

type Scheduler struct {
	Jobs      []*core.Job
	Machines  []*core.Machine
	Templates map[string]core.JobTemplate

	Planners []core.Strategy

	machineTypeRegistry map[string]core.MachineType
	jobCounter          int
	operationCounter    core.OperationID
	machineCounter      int

	engine *engine.Engine
}

func (s *Scheduler) SetLogger(l *zap.Logger) {
	s.engine.SetLogger(l)
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
	f.engine = engine.New(planners...)
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

func (f *Scheduler) Plan(problem *core.Problem) ([]engine.Report, error) {
	return f.engine.Solve(problem)
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
