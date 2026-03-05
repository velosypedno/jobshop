package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/internal/chart"
	"github.com/velosypedno/resource-allocation/internal/factory"
	"github.com/velosypedno/resource-allocation/internal/parser"
	"github.com/velosypedno/resource-allocation/internal/strategy/annealing"
	"github.com/velosypedno/resource-allocation/internal/strategy/naive"
	"github.com/velosypedno/resource-allocation/internal/strategy/rnd"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <factory_config_path> <orders_path> [optional_name]")
		os.Exit(1)
	}

	factoryConfigPath := os.Args[1]
	ordersPath := os.Args[2]
	customName := ""
	if len(os.Args) > 3 {
		customName = "_" + os.Args[3]
	}

	machinesConfig, templates, err := parser.ParseFactoryConfig(factoryConfigPath)
	if err != nil {
		fmt.Printf("Error parsing factory config: %v\n", err)
		os.Exit(1)
	}

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

	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	orders, err := parser.ParseOrders(ordersPath)
	if err != nil {
		fmt.Printf("Error parsing orders: %v\n", err)
		os.Exit(1)
	}

	results, err := f.Plan(orders, startTime)
	if err != nil {
		fmt.Printf("Error during planning: %v\n", err)
		os.Exit(1)
	}

	solutionsChart := chart.GenerateFromSolutions(results, f.Machines)

	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("plan_%s%s.html", timestamp, customName)
	fullPath := filepath.Join(outputDir, fileName)

	outputFile, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = solutionsChart.Render(outputFile)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", fullPath)
}
