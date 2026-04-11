package base

import (
	"fmt"
	"strings"
	"time"
)

type (
	OperationID int
	JobID       int
)

type OperationTemplate struct {
	Name           string
	MachineType    MachineType
	ProcessingTime time.Duration
	Children       []OperationTemplate
}

type Operation struct {
	ID OperationID

	Name            string
	MachineType     MachineType
	Duration        time.Duration
	ChildOperations []*Operation
}

type JobTemplate struct {
	Name       string
	Operations []OperationTemplate
}

type Job struct {
	ID         JobID
	Name       string
	Operations []*Operation
}

func CreateJob(id JobID, template JobTemplate, lastOperationID *OperationID) Job {
	job := Job{
		ID:         id,
		Name:       template.Name,
		Operations: []*Operation{},
	}

	for _, operation := range template.Operations {
		job.Operations = append(job.Operations, instantiateOperation(id, operation, lastOperationID))
	}

	return job
}

func instantiateOperation(jobID JobID, t OperationTemplate, lastOperationID *OperationID) *Operation {
	instance := Operation{
		ID:          OperationID(*lastOperationID),
		Name:        t.Name,
		MachineType: t.MachineType,
		Duration:    t.ProcessingTime,
	}
	*lastOperationID++

	for _, child := range t.Children {
		instance.ChildOperations = append(instance.ChildOperations, instantiateOperation(jobID, child, lastOperationID))
	}

	return &instance
}

func (j Job) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("JOB ID:   %d\n", j.ID))
	sb.WriteString(fmt.Sprintf("NAME:     %s\n", j.Name))
	sb.WriteString("==========================================\n")

	for _, op := range j.Operations {
		sb.WriteString(op.formatTree(1))
	}

	return sb.String()
}

func (o *Operation) formatTree(level int) string {
	var sb strings.Builder

	var indent string
	if level > 1 {
		indent = strings.Repeat("  │ ", level-1) + "  ├─ "
	} else {
		indent = " ├─ "
	}

	sb.WriteString(fmt.Sprintf("%s [%d] %s (%d, %v)\n",
		indent,
		o.ID,
		o.Name,
		o.MachineType,
		o.Duration,
	))

	for _, sub := range o.ChildOperations {
		sb.WriteString(sub.formatTree(level + 1))
	}
	return sb.String()
}
