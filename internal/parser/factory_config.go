package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/strategy/annealing"
	"github.com/velosypedno/resource-allocation/internal/strategy/ga"
	"github.com/velosypedno/resource-allocation/internal/strategy/naive"
	"github.com/velosypedno/resource-allocation/internal/strategy/rnd"
	"github.com/velosypedno/resource-allocation/internal/strategy/tabu"
	"go.uber.org/zap"
)

type Strategy interface {
	Plan([]*base.Job, []*base.Machine, time.Time) (*base.Solution, base.MachineTimeSlots)
	Name() string
	Description() string
	SetLogger(l *zap.Logger)
}

func ParseFactoryConfig(filePath string) ([]MachineConfig, []base.JobTemplate, []Strategy, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config FactoryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse json: %w", err)
	}

	machineTypeMap := make(map[string]base.MachineType)
	for _, m := range config.Machines {
		machineTypeMap[m.TypeName] = base.MachineType(m.TypeID)
	}

	templates := make([]base.JobTemplate, 0, len(config.JobTemplates))
	for _, j := range config.JobTemplates {
		operations, err := convertOperations(j.Operations, machineTypeMap)
		if err != nil {
			return nil, nil, nil, err
		}
		templates = append(templates, base.JobTemplate{
			Name:       j.Name,
			Operations: operations,
		})
	}

	strategies := make([]Strategy, 0, len(config.Strategies))
	for _, sDTO := range config.Strategies {
		s, err := createStrategy(sDTO)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("strategy '%s': %w", sDTO.Name, err)
		}
		strategies = append(strategies, s)
	}

	return config.Machines, templates, strategies, nil
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

func createStrategy(dto StrategyDTO) (Strategy, error) {
	switch dto.Type {
	case "ga":
		var p GAConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}
		return ga.New(p.PopulationSize, p.Generations, p.MutationRate, p.CrossoverRate, p.ElitismRatio), nil

	case "tabu":
		var p TabuConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}
		return tabu.New(p.TabuSize, p.MaxIterations, p.NeighborsCount), nil

	case "annealing_priority_based":
		var p AnnealingConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}
		annealingConfig := annealing.Config{
			InitialTemp:      p.InitialTemp,
			MinTemp:          p.MinTemp,
			Alpha:            p.Alpha,
			Iterations:       p.Iterations,
			SwapsPerMutation: p.Swaps,
		}
		return annealing.NewPriorityBased(annealingConfig), nil

	case "annealing_sequence_based":
		var p AnnealingConfigDTO
		if err := json.Unmarshal(dto.Params, &p); err != nil {
			return nil, err
		}
		annealingConfig := annealing.Config{
			InitialTemp:      p.InitialTemp,
			MinTemp:          p.MinTemp,
			Alpha:            p.Alpha,
			Iterations:       p.Iterations,
			SwapsPerMutation: p.Swaps,
		}
		return annealing.NewSequenceBased(annealingConfig), nil

	case "greedy", "naive":
		return naive.New(), nil

	case "random":
		return rnd.New(), nil

	default:
		return nil, fmt.Errorf("unknown strategy type: %s", dto.Type)
	}
}
