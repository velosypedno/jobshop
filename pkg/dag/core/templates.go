package core

import (
	"time"
)

type (
	OpTemplateID  int
	JobTemplateID int
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
	ID JobTemplateID

	Name  string
	Ops   map[OpTemplateID]OpTemplate
	Edges []OpTemplateDependencyEdge
}

func InitJobTemplate(name string) JobGraphTemplate {
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

func validateOperationsTopology(jobGraphTemplate JobGraphTemplate) error {
	// check all edge endpoint IDs exist
	for _, edge := range jobGraphTemplate.Edges {
		if _, ok := jobGraphTemplate.Ops[edge.From]; !ok {
			return &ErrUnknownOpTemplateDetail{ID: edge.From}
		}
		if _, ok := jobGraphTemplate.Ops[edge.To]; !ok {
			return &ErrUnknownOpTemplateDetail{ID: edge.To}
		}
	}

	// detect cycles via Kahn's algorithm
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
		return ErrCycleInGraph
	}

	return nil
}
