package core

import (
	"fmt"
	"time"
)

type (
	MachineType  uint
	MachineID    int
	OpID         int
	JobID        int
	OpTemplateID int
)

type OpTemplate struct {
	ID OpTemplateID

	Name           string
	MachineType    MachineType
	ProcessingTime time.Duration
}

type OpTemplateDependencyEdge struct {
	From OpTemplateID
	To   OpTemplateID
}

type JobGraphTemplate struct {
	Name  string
	Ops   map[OpTemplateID]OpTemplate
	Edges []OpTemplateDependencyEdge
}

func InitJobtemplate(name string) JobGraphTemplate {
	return JobGraphTemplate{
		Name:  name,
		Ops:   make(map[OpTemplateID]OpTemplate),
		Edges: []OpTemplateDependencyEdge{},
	}
}

func (j *JobGraphTemplate) AddOpTemplate(opTemplate OpTemplate) {
	j.Ops[opTemplate.ID] = opTemplate
}

func (j *JobGraphTemplate) AddEdge(from, to OpTemplateID) {
	j.Edges = append(j.Edges, OpTemplateDependencyEdge{from, to})
}

func ValidateTemplate(jobGraphTemplate JobGraphTemplate) error {
	// cheсk all IDs exist
	for _, edge := range jobGraphTemplate.Edges {
		if _, ok := jobGraphTemplate.Ops[edge.From]; !ok {
			return fmt.Errorf("template with id %v does not exist", edge.From)
		}
		if _, ok := jobGraphTemplate.Ops[edge.To]; !ok {
			return fmt.Errorf("template with id %v does not exist", edge.From)
		}
	}

	// check cycle exists
	inDegreeMap := make(map[OpTemplateID]int, len(jobGraphTemplate.Ops))
	for _, opTemplate := range jobGraphTemplate.Ops {
		inDegreeMap[opTemplate.ID] = 0
	}
	for _, edge := range jobGraphTemplate.Edges {
		inDegreeMap[edge.To]++
	}

	successorsMap := make(map[OpTemplateID][]OpTemplateID)
	for _, edge := range jobGraphTemplate.Edges {
		successorsMap[edge.From] = append(successorsMap[edge.From], edge.To)
	}

	sources := []OpTemplateID{}
	for opTemplateID, inDegreeValue := range inDegreeMap {
		if inDegreeValue == 0 {
			sources = append(sources, opTemplateID)
		}
	}
	unvisited := len(inDegreeMap) - len(sources)
	for {
		if len(sources) == 0 {
			break
		}

		source := sources[0]
		sources = sources[1:]
		successorsTemplatesIDs := successorsMap[source]
		for _, successorTemplateID := range successorsTemplatesIDs {
			inDegreeMap[successorTemplateID]--
			if inDegreeMap[successorTemplateID] == 0 {
				unvisited--
				sources = append(sources, successorTemplateID)
			}
		}
	}

	if unvisited != 0 {
		return fmt.Errorf("cycle presents in topology")
	}

	return nil
}

type Operation struct {
	ID OpID

	Name           string
	ProcessingTime time.Duration
	MachineType    MachineType

	Predecessors []OpID
	Successors   []OpID
}

type Job struct {
	ID JobID

	Name string
	Ops  map[OpID]*Operation
}

func CreateJob(id JobID, template JobGraphTemplate, lastOperationID *OpID) (*Job, error) {
	err := ValidateTemplate(template)
	if err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	job := &Job{
		ID:   id,
		Name: template.Name,
		Ops:  make(map[OpID]*Operation, len(template.Ops)),
	}

	opTempalateIDToActualOpID := make(map[OpTemplateID]OpID, len(template.Ops))

	for _, operationTemplate := range template.Ops {
		*lastOperationID++

		opTempalateIDToActualOpID[operationTemplate.ID] = *lastOperationID
		job.Ops[*lastOperationID] = &Operation{
			ID:             *lastOperationID,
			Name:           operationTemplate.Name,
			ProcessingTime: operationTemplate.ProcessingTime,
			MachineType:    operationTemplate.MachineType,
		}
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

	return job, nil
}
