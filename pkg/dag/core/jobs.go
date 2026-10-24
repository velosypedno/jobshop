package core

import (
	"time"
)

type (
	MachineType uint
	MachineID   int
	OpID        int
	JobID       int
)

// readonly
type Operation struct {
	ID OpID

	Name           string
	ProcessingTime time.Duration
	MachineType    MachineType

	Predecessors []OpID
	Successors   []OpID
}

// readonly
type Job struct {
	ID JobID

	Name string
	Ops  map[OpID]*Operation
}

// template shoud be check for correctness before using
func createJob(id JobID, template JobGraphTemplate, lastOperationID *OpID) *Job {
	job := &Job{
		ID:   id,
		Name: template.Name,
		Ops:  make(map[OpID]*Operation, len(template.Ops)),
	}

	opTempalateIDToActualOpID := make(map[OpTemplateID]OpID, len(template.Ops))

	for _, operationTemplate := range template.Ops {

		opTempalateIDToActualOpID[operationTemplate.ID] = *lastOperationID
		job.Ops[*lastOperationID] = &Operation{
			ID:             *lastOperationID,
			Name:           operationTemplate.Name,
			ProcessingTime: operationTemplate.ProcessingTime,
			MachineType:    operationTemplate.MachineType,
		}
		*lastOperationID++
	}

	// get toplogy map
	opTemlpateIDsEdgeMap := make(map[OpTemplateID][]OpTemplateID)
	for _, edge := range template.Edges {
		opTemlpateIDsEdgeMap[edge.From] = append(opTemlpateIDsEdgeMap[edge.From], edge.To)
	}

	// enrcih with actual successors and predecessors ids
	for predecessorTemplateID, successorsTemplatesIDs := range opTemlpateIDsEdgeMap {
		predecessorOpId := opTempalateIDToActualOpID[predecessorTemplateID]
		predecessorOp := job.Ops[predecessorOpId]
		for _, successorTemplateID := range successorsTemplatesIDs {
			successorOpId := opTempalateIDToActualOpID[successorTemplateID]
			predecessorOp.Successors = append(predecessorOp.Successors, successorOpId)
			successorOp := job.Ops[successorOpId]
			successorOp.Predecessors = append(successorOp.Predecessors, predecessorOpId)
		}
	}

	return job
}
