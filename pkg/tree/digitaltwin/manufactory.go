package digitaltwin

import (
	"fmt"
	"time"

	"github.com/velosypedno/jobshop/pkg/tree/core"
)

type Manufactory struct {
	Machines  []*core.Machine
	Templates map[string]core.JobTemplate
}

func New(machineConfigs []core.MachineConfigEntry, templates []core.JobTemplate) *Manufactory {
	f := &Manufactory{}
	f.Templates = make(map[string]core.JobTemplate)

	for _, t := range templates {
		f.Templates[t.Name] = t
	}

	var machineCounter int
	for _, mConf := range machineConfigs {
		mType := mConf.Type

		for i := 0; i < mConf.Count; i++ {
			machineCounter++
			m := core.NewMachine(core.MachineID(machineCounter), mType, mConf.Name)
			m.Name = mConf.Name
			f.Machines = append(f.Machines, &m)
		}
	}
	return f
}

func (f *Manufactory) NewProblem(orders []core.OrderEntry, startTime time.Time) *core.Problem {
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

func (f *Manufactory) createJobsFromOrders(orders []core.OrderEntry) ([]*core.Job, error) {
	var jobs []*core.Job
	jobIDCounter := core.JobID(0)

	operationCounter := new(core.OperationID)
	for _, order := range orders {
		template, ok := f.Templates[order.Name]
		if !ok {
			return nil, fmt.Errorf("template '%s' not found for order", order.Name)
		}

		for i := 0; i < order.Amount; i++ {
			jobIDCounter++
			newJob := core.CreateJob(jobIDCounter, template, operationCounter)
			jobs = append(jobs, &newJob)
		}
	}
	return jobs, nil
}
