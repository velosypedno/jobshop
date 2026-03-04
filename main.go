package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/chart"
	"github.com/velosypedno/resource-allocation/factory"
	"github.com/velosypedno/resource-allocation/parser"
	"github.com/velosypedno/resource-allocation/strategy/annealing"
	"github.com/velosypedno/resource-allocation/strategy/naive"
	"github.com/velosypedno/resource-allocation/strategy/rnd"
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
	annealingStrategy := annealing.New(100, 0.1, 0.99, 100, 15)
	f.SetPlanners(&rnd.Strategy{}, &naive.Strategy{}, annealingStrategy)

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
