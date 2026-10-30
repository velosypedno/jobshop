package core

import (
	"encoding/json"
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

// allow machines with different names, but equal types
type MachineConfigEntry struct {
	Type  MachineType
	Name  string
	Count int
}

type FactoryConfig struct {
	MachineConfigs []MachineConfigEntry
	JobTemplates   []JobGraphTemplate
}

type factoryConfigDTO struct {
	Machines     []machineConfigDTO `json:"machines"`
	JobTemplates []jobTemplateDTO   `json:"job_templates"`
}

type machineConfigDTO struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
	Count    int    `json:"count"`
}

type jobTemplateDTO struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	Operations []operationTemplateDTO `json:"operations"`
	Edges      []edgeDTO              `json:"edges"`
}

type edgeDTO struct {
	FromOpID int `json:"from_op_id"`
	ToOpID   int `json:"to_op_id"`
}

type operationTemplateDTO struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	MachineType    int    `json:"machine_type"`
	ProcessingTime string `json:"processing_time"`
}

func (fc *FactoryConfig) ToJSON() ([]byte, error) {
	dto := factoryConfigDTO{
		Machines:     make([]machineConfigDTO, len(fc.MachineConfigs)),
		JobTemplates: make([]jobTemplateDTO, len(fc.JobTemplates)),
	}

	for i, m := range fc.MachineConfigs {
		dto.Machines[i] = machineConfigDTO{
			TypeID:   int(m.Type),
			TypeName: m.Name,
			Count:    m.Count,
		}
	}

	for i, jt := range fc.JobTemplates {
		ops := make([]operationTemplateDTO, 0, len(jt.Ops))
		for _, op := range jt.Ops {
			ops = append(ops, operationTemplateDTO{
				ID:             int(op.ID),
				Name:           op.Name,
				MachineType:    int(op.MachineType),
				ProcessingTime: op.ProcessingTime.String(),
			})
		}

		edges := make([]edgeDTO, len(jt.Edges))
		for j, e := range jt.Edges {
			edges[j] = edgeDTO{
				FromOpID: int(e.From),
				ToOpID:   int(e.To),
			}
		}

		dto.JobTemplates[i] = jobTemplateDTO{
			ID:         int(jt.ID),
			Name:       jt.Name,
			Operations: ops,
			Edges:      edges,
		}
	}

	return json.Marshal(dto)
}

func FromJSON(data []byte) (FactoryConfig, error) {
	var dto factoryConfigDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return FactoryConfig{}, err
	}

	machineConfigs := make([]MachineConfigEntry, len(dto.Machines))
	for i, m := range dto.Machines {
		machineConfigs[i] = MachineConfigEntry{
			Type:  MachineType(m.TypeID),
			Name:  m.TypeName,
			Count: m.Count,
		}
	}

	jobTemplates := make([]JobGraphTemplate, len(dto.JobTemplates))
	for i, jtDTO := range dto.JobTemplates {
		jt := JobGraphTemplate{
			ID:    JobTemplateID(jtDTO.ID),
			Name:  jtDTO.Name,
			Ops:   make(map[OpTemplateID]OpTemplate),
			Edges: make([]OpTemplateDependencyEdge, len(jtDTO.Edges)),
		}

		for _, opDTO := range jtDTO.Operations {
			dur, err := time.ParseDuration(opDTO.ProcessingTime)
			if err != nil {
				return FactoryConfig{}, err
			}
			opID := OpTemplateID(opDTO.ID)
			jt.Ops[opID] = OpTemplate{
				ID:             opID,
				Name:           opDTO.Name,
				MachineType:    MachineType(opDTO.MachineType),
				ProcessingTime: dur,
			}
		}

		for j, eDTO := range jtDTO.Edges {
			jt.Edges[j] = OpTemplateDependencyEdge{
				From: OpTemplateID(eDTO.FromOpID),
				To:   OpTemplateID(eDTO.ToOpID),
			}
		}

		jobTemplates[i] = jt
	}

	return FactoryConfig{
		MachineConfigs: machineConfigs,
		JobTemplates:   jobTemplates,
	}, nil
}
