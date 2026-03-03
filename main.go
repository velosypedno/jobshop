package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/velosypedno/resource-allocation/chart"
	"github.com/velosypedno/resource-allocation/factory"
	"github.com/velosypedno/resource-allocation/parser"
	"github.com/velosypedno/resource-allocation/strategy/naive"
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
		customName = os.Args[3] + "_"
	}

	machinesConfig, templates, err := parser.ParseFactoryConfig(factoryConfigPath)
	if err != nil {
		fmt.Printf("Error parsing factory config: %v\n", err)
		os.Exit(1)
	}

	f := &factory.Factory{}
	f.Configure(machinesConfig, templates)
	f.SetPlanner(&naive.Strategy{})

	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local)

	orders, err := parser.ParseOrders(ordersPath)
	if err != nil {
		fmt.Printf("Error parsing orders: %v\n", err)
		os.Exit(1)
	}

	solution, metaInfo, err := f.Plan(orders, startTime)
	if err != nil {
		fmt.Printf("Error during planning: %v\n", err)
		os.Exit(1)
	}

	solutionChart := chart.GenerateFromSolution(solution, f.Machines, metaInfo)

	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("%splan_%s.html", customName, timestamp)
	fullPath := filepath.Join(outputDir, fileName)

	outputFile, err := os.Create(fullPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = solutionChart.Render(outputFile)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", fullPath)
}
