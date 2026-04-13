package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/jobshop/internal/core"
	"github.com/velosypedno/jobshop/internal/strategy/annealing"
	"github.com/velosypedno/jobshop/internal/strategy/ga"
	"github.com/velosypedno/jobshop/internal/strategy/naive"
	"github.com/velosypedno/jobshop/internal/strategy/tabu"
)

func ParseFactoryConfig(filePath string) ([]MachineConfig, []core.JobTemplate, []core.Strategy, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config FactoryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse json: %w", err)
	}

	machineTypeMap := make(map[string]core.MachineType)
	for _, m := range config.Machines {
		machineTypeMap[m.TypeName] = core.MachineType(m.TypeID)
	}

	templates := make([]core.JobTemplate, 0, len(config.JobTemplates))
	for _, j := range config.JobTemplates {
		operations, err := convertOperations(j.Operations, machineTypeMap)
		if err != nil {
			return nil, nil, nil, err
		}
		templates = append(templates, core.JobTemplate{
			Name:       j.Name,
			Operations: operations,
		})
	}

	strategies := make([]core.Strategy, 0, len(config.Strategies))
	for _, sDTO := range config.Strategies {
		s, err := createStrategy(sDTO)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("strategy '%s': %w", sDTO.Name, err)
		}
		strategies = append(strategies, s)
	}

	return config.Machines, templates, strategies, nil
}

func convertOperations(dtos []OperationTemplateDTO, machineTypes map[string]core.MachineType) ([]core.OperationTemplate, error) {
	if dtos == nil {
		return nil, nil
	}

	res := make([]core.OperationTemplate, len(dtos))
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

		res[i] = core.OperationTemplate{
			Name:           d.Name,
			MachineType:    mType,
			ProcessingTime: duration,
			Children:       children,
		}
	}
	return res, nil
}

func createStrategy(dto StrategyDTO) (core.Strategy, error) {
	switch dto.Type {
	case "ga":
		var p GAConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}
		return ga.New(p.PopulationSize, p.Generations, p.MutationRate, p.CrossoverRate, p.ElitismRatio, dto.Name), nil

	case "tabu":
		var p TabuConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}
		return tabu.New(p.TabuSize, p.MaxIterations, p.NeighborsCount, dto.Name), nil

	case "annealing_priority_based":
		var p AnnealingConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}

		return annealing.New(p.InitialTemp, p.MinTemp, p.Alpha, p.Iterations, p.Swaps, dto.Name), nil

	case "greedy", "naive":
		return naive.New(dto.Name), nil

	default:
		return nil, fmt.Errorf("unknown strategy type: %s", dto.Type)
	}
}
