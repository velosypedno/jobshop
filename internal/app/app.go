package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/internal/base"
	"github.com/velosypedno/resource-allocation/internal/chart"
	"github.com/velosypedno/resource-allocation/internal/factory"
	"github.com/velosypedno/resource-allocation/internal/parser"
	"github.com/velosypedno/resource-allocation/internal/strategy/annealing"
	"github.com/velosypedno/resource-allocation/internal/strategy/naive"
	"github.com/velosypedno/resource-allocation/internal/strategy/rnd"
)

type App struct {
	Factory *factory.Factory
}

func New(machinesConfig []parser.MachineConfig, templates []base.JobTemplate) *App {
	f := &factory.Factory{}
	f.Configure(machinesConfig, templates)

	annealingConfig := annealing.Config{
		InitialTemp:      100,
		MinTemp:          0.1,
		Alpha:            0.99,
		Iterations:       100,
		SwapsPerMutation: 15,
	}
	sequenceBasedAnnealing := annealing.NewSequenceBased(annealingConfig)
	priorityBasedAnnealing := annealing.NewPriorityBased(annealingConfig)
	randomStrategy := rnd.New()
	naiveStrategy := naive.New()

	f.SetPlanners(
		randomStrategy,
		naiveStrategy,
		sequenceBasedAnnealing,
		priorityBasedAnnealing,
	)

	return &App{
		Factory: f,
	}
}

func (a *App) Run(startTime time.Time, orders []parser.OrderDTO, customName string) error {
	results, err := a.Factory.Plan(orders, startTime)
	if err != nil {
		return fmt.Errorf("during planning: %v", err)

	}

	solutionsChart := chart.GenerateFromSolutions(results, a.Factory.Machines)

	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating directory: %v", err)

	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("plan_%s%s.html", timestamp, customName)
	fullPath := filepath.Join(outputDir, fileName)

	outputFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer outputFile.Close()

	err = solutionsChart.Render(outputFile)
	if err != nil {
		return fmt.Errorf("error rendering chart: %v", err)
	}

	fmt.Printf("Successfully generated %s\n", fullPath)
	return nil
}
