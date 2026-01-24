package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/resource-allocation/base"
)

func ParseFactoryConfig(filePath string) ([]MachineConfig, []base.JobTemplate, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config FactoryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, nil, fmt.Errorf("failed to parse json: %w", err)
	}

	machineTypeMap := make(map[string]base.MachineType)
	for _, m := range config.Machines {
		machineTypeMap[m.TypeName] = base.MachineType(m.TypeID)
	}

	templates := make([]base.JobTemplate, 0, len(config.JobTemplates))
	for _, j := range config.JobTemplates {
		operations, err := convertOperations(j.Operations, machineTypeMap)
		if err != nil {
			return nil, nil, err
		}
		templates = append(templates, base.JobTemplate{
			Name:       j.Name,
			Operations: operations,
		})
	}

	return config.Machines, templates, nil
}

func convertOperations(dtos []OperationTemplateDTO, machineTypes map[string]base.MachineType) ([]base.OperationTemplate, error) {
	if dtos == nil {
		return nil, nil
	}

	res := make([]base.OperationTemplate, len(dtos))
	for i, d := range dtos {
		duration, err := time.ParseDuration(d.ProcessingTime)
		if err != nil {
			return nil, fmt.Errorf("operation '%s': invalid duration format '%s' (example: '60m', '1h5s')", d.Name, d.ProcessingTime)
		}

		mType, ok := machineTypes[d.MachineType]
		if !ok {
			return nil, fmt.Errorf("operation '%s': machine type '%s' is not defined in the factory configuration", d.Name, d.MachineType)
		}

		children, err := convertOperations(d.Children, machineTypes)
		if err != nil {
			return nil, fmt.Errorf("in child of '%s' -> %w", d.Name, err)
		}

		res[i] = base.OperationTemplate{
			Name:           d.Name,
			MachineType:    mType,
			ProcessingTime: duration,
			Children:       children,
		}
	}
	return res, nil
}
