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
	"github.com/velosypedno/resource-allocation/internal/strategy/ga"
	"github.com/velosypedno/resource-allocation/internal/strategy/naive"
	"github.com/velosypedno/resource-allocation/internal/strategy/rnd"
	"github.com/velosypedno/resource-allocation/internal/strategy/tabu"
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
	gaStrategy := ga.New(16, 100, 0.05, 0.7, 0.125)
	anotherGAStrategy := ga.New(64, 1000, 0.05, 0.7, 0.125)
	morePopulationGAStrategy := ga.New(128, 1000, 0.05, 0.7, 0.125)
	giantPopulationGAStrategy := ga.New(256, 2000, 0.03, 0.65, 0.125)
	tabuStrategy := tabu.New(15, 500, 30)

	f.SetPlanners(
		randomStrategy,
		naiveStrategy,
		sequenceBasedAnnealing,
		priorityBasedAnnealing,
		tabuStrategy,
		gaStrategy,
		anotherGAStrategy,
		morePopulationGAStrategy,
		giantPopulationGAStrategy,
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
