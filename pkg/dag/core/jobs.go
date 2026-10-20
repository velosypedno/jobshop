package core

import (
	"time"
)

type (
	MachineType         uint
	MachineID           int
	OperationID         int
	JobID               int
	OperationTemplateID int
)

type OperationTemplate struct {
	ID OperationTemplateID

	Name           string
	MachineType    MachineType
	ProcessingTime time.Duration
}

type OperationTemplateDependencyEdge struct {
	From OperationTemplateID
	To   OperationTemplateID
}

type JobGraphTemplate struct {
	Name       string
	Operations map[OperationTemplateID]OperationTemplate
	Edges      []OperationTemplateDependencyEdge
}

type Operation struct {
	ID OperationID

	Name           string
	ProcessingTime time.Duration
	MachineType    MachineType

	Predecessors []OperationID
	Successors   []OperationID
}

type Job struct {
	ID JobID

	Name       string
	Operations map[OperationID]*Operation
}

func CreateJob(id JobID, template JobGraphTemplate, lastOperationID *OperationID) *Job {
	job := &Job{
		ID:         id,
		Name:       template.Name,
		Operations: make(map[OperationID]*Operation, len(template.Operations)),
	}

	operationTempalateIDToActualOperationID := make(map[OperationTemplateID]OperationID, len(template.Operations))

	for _, operationTemplate := range template.Operations {
		*lastOperationID++

		operationTempalateIDToActualOperationID[operationTemplate.ID] = *lastOperationID
		job.Operations[*lastOperationID] = &Operation{
			ID:             *lastOperationID,
			Name:           operationTemplate.Name,
			ProcessingTime: operationTemplate.ProcessingTime,
			MachineType:    operationTemplate.MachineType,
		}
	}

	// get toplogy map
	opTemlpateIdsEdgeMap := make(map[OperationTemplateID][]OperationTemplateID)
	for _, edge := range template.Edges {
		opTemlpateIdsEdgeMap[edge.From] = append(opTemlpateIdsEdgeMap[edge.From], edge.To)
	}

	// enrcih with actual successors and predecessors ids
	for predecessorTemplateID, successorsTemplatesIDs := range opTemlpateIdsEdgeMap {
		predecessorOperationId := operationTempalateIDToActualOperationID[predecessorTemplateID]
		predecessorOperation := job.Operations[predecessorOperationId]
		for _, successorTemplateID := range successorsTemplatesIDs {
			successorOperationId := operationTempalateIDToActualOperationID[successorTemplateID]
			predecessorOperation.Successors = append(predecessorOperation.Successors, successorOperationId)
			successorOperation := job.Operations[successorOperationId]
			successorOperation.Predecessors = append(successorOperation.Predecessors, predecessorOperationId)
		}
	}

	return job
}
