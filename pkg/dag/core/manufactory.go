package core

import (
	"fmt"
)

// readonly
type Problem struct {
	Jobs     []*Job
	Machines []*Machine
}

type Manufactory struct {
	Machines  []*Machine
	Templates map[JobTemplateID]JobGraphTemplate
}

// No checks or constraints for machineConfigs are needed.
//
// 1. We do not check MachineType for uniqueness.
//
// 2. So allowed entries with the same MachineType, but different names.
//
// 3. Even entirely duplicated MachineConfigEntry are allowed, they will not ruin invariants, and will not cause errors
func NewManufactory(machineConfigs []MachineConfigEntry, jobTemplates []JobGraphTemplate) (*Manufactory, error) {
	f := &Manufactory{}

	machineTypeSet := make(map[MachineType]struct{}, len(machineConfigs))

	var machineCounter int
	for _, mConf := range machineConfigs {
		mType := mConf.Type
		machineTypeSet[mType] = struct{}{}

		for i := 0; i < mConf.Count; i++ {
			machineCounter++
			m := NewMachine(MachineID(machineCounter), mType, mConf.Name)
			f.Machines = append(f.Machines, &m)
		}
	}

	for _, jobTemplate := range jobTemplates {
		for _, opTemplate := range jobTemplate.Ops {
			if _, ok := machineTypeSet[opTemplate.MachineType]; !ok {
				return nil, &ErrMissingMachineTypeDetail{MachineType: opTemplate.MachineType}
			}
		}
	}

	f.Templates = make(map[JobTemplateID]JobGraphTemplate)
	for _, jobTemplate := range jobTemplates {
		err := validateOperationsTopology(jobTemplate)
		if err != nil {
			return nil, fmt.Errorf("manufactory initialization: %w", err)
		}
		f.Templates[jobTemplate.ID] = jobTemplate
	}

	return f, nil
}

type OrderEntry struct {
	JobTemplateID
	Amount int
}

func (f *Manufactory) NewProblem(orders []OrderEntry) (*Problem, error) {
	for _, order := range orders {
		if _, ok := f.Templates[order.JobTemplateID]; !ok {
			return nil, &ErrUnknownJobTemplateDetail{ID: order.JobTemplateID}
		}
	}

	jobs := f.createJobsFromOrders(orders)

	problem := Problem{
		Jobs:     jobs,
		Machines: f.Machines,
	}
	return &problem, nil
}

func (f *Manufactory) createJobsFromOrders(orders []OrderEntry) []*Job {
	var jobs []*Job
	jobIDCounter := JobID(0)

	operationCounter := new(OpID)
	for _, order := range orders {
		template := f.Templates[order.JobTemplateID]

		for i := 0; i < order.Amount; i++ {
			jobIDCounter++
			newJob := createJob(jobIDCounter, template, operationCounter)
			jobs = append(jobs, newJob)
		}
	}
	return jobs
}
