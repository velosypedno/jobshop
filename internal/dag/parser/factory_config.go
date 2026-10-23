package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/jobshop/pkg/dag/core"
)

func ParseFactoryConfig(filePath string) ([]core.MachineConfigEntry, []core.JobGraphTemplate, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config factoryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, nil, fmt.Errorf("failed to parse json: %w", err)
	}

	templates := make([]core.JobGraphTemplate, 0, len(config.JobTemplates))
	for _, j := range config.JobTemplates {
		ops, err := convertOps(j.Operations)
		if err != nil {
			return nil, nil, err
		}

		edges := make([]core.OpTemplateDependencyEdge, len(j.Edges))
		for i, edgeDTO := range j.Edges {
			edges[i] = core.OpTemplateDependencyEdge{
				From: core.OpTemplateID(edgeDTO.FromOpID),
				To:   core.OpTemplateID(edgeDTO.ToOpID),
			}
		}

		templates = append(templates, core.JobGraphTemplate{
			ID: core.JobTemplateID(j.ID),

			Name:  j.Name,
			Ops:   ops,
			Edges: edges,
		})
	}

	machineConfigEntries := make([]core.MachineConfigEntry, 0, len(config.Machines))
	for _, machineConfigDTO := range config.Machines {
		var mConfigEntry core.MachineConfigEntry
		mConfigEntry.Count = machineConfigDTO.Count
		mConfigEntry.Name = machineConfigDTO.TypeName
		mConfigEntry.Type = core.MachineType(machineConfigDTO.TypeID)
		machineConfigEntries = append(machineConfigEntries, mConfigEntry)
	}

	return machineConfigEntries, templates, nil
}

func convertOps(dtos []operationTemplateDTO) (map[core.OpTemplateID]core.OpTemplate, error) {
	if dtos == nil {
		return nil, nil
	}

	res := make(map[core.OpTemplateID]core.OpTemplate, len(dtos))
	for _, d := range dtos {
		duration, err := time.ParseDuration(d.ProcessingTime)
		if err != nil {
			return nil, fmt.Errorf("operation '%s': invalid duration format '%s' (example: '60m', '1h5s')", d.Name, d.ProcessingTime)
		}

		res[core.OpTemplateID(d.ID)] = core.OpTemplate{
			ID: core.OpTemplateID(d.ID),

			Name:           d.Name,
			MachineType:    core.MachineType(d.MachineType),
			ProcessingTime: duration,
		}
	}
	return res, nil
}
