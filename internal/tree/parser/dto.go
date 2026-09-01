package parser

import "encoding/json"

type FactoryConfig struct {
	Machines     []MachineConfig  `json:"machines"`
	JobTemplates []JobTemplateDTO `json:"job_templates"`
	Strategies   []StrategyDTO    `json:"strategies"`
}

type StrategyDTO struct {
	Type   string          `json:"type"`
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

type GAConfigDTO struct {
	PopulationSize int     `json:"population_size"`
	Generations    int     `json:"generations"`
	MutationRate   float64 `json:"mutation_rate"`
	CrossoverRate  float64 `json:"crossover_rate"`
	ElitismRatio   float64 `json:"elitism_ratio"`
}

type TabuConfigDTO struct {
	TabuSize       int `json:"tabu_size"`
	MaxIterations  int `json:"max_iterations"`
	NeighborsCount int `json:"neighbors_count"`
}

type AnnealingConfigDTO struct {
	InitialTemp float64 `json:"initial_temp"`
	MinTemp     float64 `json:"min_temp"`
	Alpha       float64 `json:"alpha"`
	Iterations  int     `json:"iterations"`
	Swaps       int     `json:"swaps"`
}

type MachineConfig struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
	Count    int    `json:"count"`
}

type JobTemplateDTO struct {
	Name       string                 `json:"name"`
	Operations []OperationTemplateDTO `json:"operations"`
}

type OperationTemplateDTO struct {
	Name           string                 `json:"name"`
	MachineType    string                 `json:"machine_type"`
	ProcessingTime string                 `json:"processing_time"`
	Children       []OperationTemplateDTO `json:"children,omitempty"`
}

type OrderConfig struct {
	Orders []OrderDTO `json:"orders"`
}

type OrderDTO struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type StrategiesConfig struct {
	Strategies []StrategyDTO `json:"strategies"`
}
