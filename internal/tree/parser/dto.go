package parser

import "encoding/json"

type factoryConfig struct {
	Machines     []machineConfigDTO `json:"machines"`
	JobTemplates []jobTemplateDTO   `json:"job_templates"`
	Strategies   []strategyDTO      `json:"strategies"`
}

type machineConfigDTO struct {
	TypeID   int    `json:"type_id"`
	TypeName string `json:"type_name"`
	Count    int    `json:"count"`
}

type (
	strategyDTO struct {
		Type   string          `json:"type"`
		Name   string          `json:"name"`
		Params json.RawMessage `json:"params"`
	}

	gaConfigDTO struct {
		PopulationSize int     `json:"population_size"`
		Generations    int     `json:"generations"`
		MutationRate   float64 `json:"mutation_rate"`
		CrossoverRate  float64 `json:"crossover_rate"`
		ElitismRatio   float64 `json:"elitism_ratio"`
	}

	tabuConfigDTO struct {
		TabuSize       int `json:"tabu_size"`
		MaxIterations  int `json:"max_iterations"`
		NeighborsCount int `json:"neighbors_count"`
	}

	annealingConfigDTO struct {
		InitialTemp float64 `json:"initial_temp"`
		MinTemp     float64 `json:"min_temp"`
		Alpha       float64 `json:"alpha"`
		Iterations  int     `json:"iterations"`
		Swaps       int     `json:"swaps"`
	}
)

type (
	jobTemplateDTO struct {
		Name       string                 `json:"name"`
		Operations []operationTemplateDTO `json:"operations"`
	}

	operationTemplateDTO struct {
		Name           string                 `json:"name"`
		MachineType    string                 `json:"machine_type"`
		ProcessingTime string                 `json:"processing_time"`
		Children       []operationTemplateDTO `json:"children,omitempty"`
	}
)
